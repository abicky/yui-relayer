package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	name = "github.com/hyperledger-labs/yui-relayer"
)

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	logger = otelslog.NewLogger(name)

	supportedExporters = []string{
		"otlp",
		"console",
	}
)

func getExporters(envName string) ([]string, error) {
	if v := os.Getenv(envName); v != "" {
		exporters := strings.Split(v, ",")
		for _, exporter := range exporters {
			if !slices.Contains(supportedExporters, exporter) {
				return nil, fmt.Errorf("unsupported exporter: %q", exporter)
			}
		}
		return exporters, nil
	}
	return nil, nil
}

// SetupOTelSDK bootstraps the OpenTelemetry pipeline using the environment variables
// described on https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#exporter-selection.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, prometheusAddr string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// TODO: Should we support propagator?
	// prop := newPropagator()
	// otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTracerProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx, prometheusAddr)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	err = initializeMetrics()
	if err != nil {
		handleErr(err)
		return
	}

	loggerProvider, err := newLoggerProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context) (*trace.TracerProvider, error) {
	exporters, err := getExporters("OTEL_TRACES_EXPORTER")
	if err != nil {
		return nil, err
	}
	if len(exporters) == 0 {
		return nil, nil
	}

	exps := make([]trace.SpanExporter, 0)
	if slices.Contains(exporters, "otlp") {
		exp, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	if slices.Contains(exporters, "console") {
		exp, err := stdouttrace.New()
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	opts := make([]trace.TracerProviderOption, 0, len(exps))
	for _, exp := range exps {
		opt := trace.WithBatcher(exp)
		opts = append(opts, opt)
	}

	return trace.NewTracerProvider(opts...), nil
}

func newMeterProvider(ctx context.Context, prometheusAddr string) (*metric.MeterProvider, error) {
	exporters, err := getExporters("OTEL_METRICS_EXPORTER")
	if err != nil {
		return nil, err
	}
	if len(exporters) == 0 && prometheusAddr == "" {
		return nil, nil
	}

	readers := make([]metric.Reader, 0)

	if slices.Contains(exporters, "otlp") {
		exp, err := otlpmetricgrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		readers = append(readers, metric.NewPeriodicReader(exp))
	}

	if slices.Contains(exporters, "console") {
		exp, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		readers = append(readers, metric.NewPeriodicReader(exp))
	}

	if slices.Contains(exporters, "prometheus") || prometheusAddr != "" {
		// TODO: Get prometheusAddr from env
		exporter, err := NewPrometheusExporter(prometheusAddr)
		if err != nil {
			return nil, err
		}
		readers = append(readers, exporter)
	}

	opts := make([]metric.Option, 0, len(readers))
	for _, reader := range readers {
		opt := metric.WithReader(reader)
		opts = append(opts, opt)
	}

	return metric.NewMeterProvider(opts...), nil
}

func newLoggerProvider(ctx context.Context) (*log.LoggerProvider, error) {
	exporters, err := getExporters("OTEL_LOGS_EXPORTER")
	if err != nil {
		return nil, err
	}
	if len(exporters) == 0 {
		return nil, nil
	}

	exps := make([]log.Exporter, 0)

	if slices.Contains(exporters, "otlp") {
		exp, err := otlploggrpc.New(ctx)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	if slices.Contains(exporters, "console") {
		exp, err := stdoutlog.New()
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}

	opts := make([]log.LoggerProviderOption, 0, len(exps))
	for _, exp := range exps {
		opt := log.WithProcessor(log.NewBatchProcessor(exp))
		opts = append(opts, opt)
	}

	return log.NewLoggerProvider(opts...), nil
}
