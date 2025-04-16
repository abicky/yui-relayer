package core

import (
	"context"
	"fmt"
	"reflect"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer = otel.Tracer("github.com/hyperledger-labs/yui-relayer/core")
)

type chainTracerBridge struct {
	Chain
}

type proverTracerBridge struct {
	Prover
}

func (b *chainTracerBridge) SendMsgs(ctx context.Context, msgs []sdk.Msg) ([]MsgID, error) {
	ctx, span := tracer.Start(ctx, "Chain.SendMsgs",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	ids, err := b.Chain.SendMsgs(ctx, msgs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return ids, err
}

func (b *chainTracerBridge) GetMsgResult(ctx context.Context, id MsgID) (MsgResult, error) {
	ctx, span := tracer.Start(ctx, "Chain.GetMsgResult",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	result, err := b.Chain.GetMsgResult(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

func (b *chainTracerBridge) LatestHeight(ctx context.Context) (ibcexported.Height, error) {
	ctx, span := tracer.Start(ctx, "Chain.LatestHeight",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	height, err := b.Chain.LatestHeight(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return height, err
}

func (b *chainTracerBridge) Timestamp(ctx context.Context, height ibcexported.Height) (time.Time, error) {
	ctx, span := tracer.Start(ctx, "Chain.Timestamp",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	t, err := b.Chain.Timestamp(ctx, height)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return t, err
}

func (b *chainTracerBridge) QueryClientConsensusState(ctx QueryContext, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryClientConsensusState",
		WithClientAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryClientConsensusState(ctx, dstClientConsHeight)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryClientState(ctx QueryContext) (*clienttypes.QueryClientStateResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryClientState",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryClientState(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryConnection(ctx QueryContext, connectionID string) (*conntypes.QueryConnectionResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryConnection",
		WithClientAttributes(b),
		trace.WithAttributes(AttributeKeyConnectionID.String(connectionID)),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryConnection(ctx, connectionID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryChannel(ctx QueryContext) (*chantypes.QueryChannelResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryChannel",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryChannel(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryUnreceivedPackets(ctx QueryContext, seqs []uint64) ([]uint64, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryUnreceivedPackets",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	packets, err := b.Chain.QueryUnreceivedPackets(ctx, seqs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return packets, err
}

func (b *chainTracerBridge) QueryUnfinalizedRelayPackets(ctx QueryContext, counterparty LightClientICS04Querier) (PacketInfoList, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryUnfinalizedRelayPackets",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	list, err := b.Chain.QueryUnfinalizedRelayPackets(ctx, counterparty)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return list, err
}

func (b *chainTracerBridge) QueryUnreceivedAcknowledgements(ctx QueryContext, seqs []uint64) ([]uint64, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryUnreceivedAcknowledgements",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	acks, err := b.Chain.QueryUnreceivedAcknowledgements(ctx, seqs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return acks, err
}

func (b *chainTracerBridge) QueryUnfinalizedRelayAcknowledgements(ctx QueryContext, counterparty LightClientICS04Querier) (PacketInfoList, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryUnfinalizedRelayAcknowledgements",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	list, err := b.Chain.QueryUnfinalizedRelayAcknowledgements(ctx, counterparty)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return list, err
}

func (b *chainTracerBridge) QueryChannelUpgrade(ctx QueryContext) (*chantypes.QueryUpgradeResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryChannelUpgrade",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryChannelUpgrade(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryChannelUpgradeError(ctx QueryContext) (*chantypes.QueryUpgradeErrorResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryChannelUpgradeError",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryChannelUpgradeError(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *chainTracerBridge) QueryCanTransitionToFlushComplete(ctx QueryContext) (bool, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryCanTransitionToFlushComplete",
		WithChannelAttributes(b),
		withPackage(b.Chain),
	)
	defer span.End()

	ok, err := b.Chain.QueryCanTransitionToFlushComplete(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return ok, err
}

func (b *chainTracerBridge) QueryBalance(ctx QueryContext, address sdk.AccAddress) (sdk.Coins, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryBalance",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	coins, err := b.Chain.QueryBalance(ctx, address)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return coins, err
}

func (b *chainTracerBridge) QueryDenomTraces(ctx QueryContext, offset, limit uint64) (*transfertypes.QueryDenomTracesResponse, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Chain.QueryDenomTraces",
		WithChainAttributes(b.ChainID()),
		withPackage(b.Chain),
	)
	defer span.End()

	resp, err := b.Chain.QueryDenomTraces(ctx, offset, limit)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return resp, err
}

func (b *proverTracerBridge) GetLatestFinalizedHeader(ctx context.Context) (Header, error) {
	ctx, span := tracer.Start(ctx, "Prover.GetLatestFinalizedHeader",
		withPackage(b.Prover),
	)
	defer span.End()

	header, err := b.Prover.GetLatestFinalizedHeader(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return header, err
}

func (b *proverTracerBridge) CreateInitialLightClientState(ctx context.Context, height ibcexported.Height) (ibcexported.ClientState, ibcexported.ConsensusState, error) {
	ctx, span := tracer.Start(ctx, "Prover.CreateInitialLightClientState",
		withPackage(b.Prover),
	)
	defer span.End()

	clientState, consensusState, err := b.Prover.CreateInitialLightClientState(ctx, height)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return clientState, consensusState, err
}

func (b *proverTracerBridge) SetupHeadersForUpdate(ctx context.Context, counterparty FinalityAwareChain, latestFinalizedHeader Header) ([]Header, error) {
	ctx, span := tracer.Start(ctx, "Prover.SetupHeadersForUpdate",
		trace.WithAttributes(attribute.String("counterparty_chain_id", counterparty.ChainID())),
		withPackage(b.Prover),
	)
	defer span.End()

	headers, err := b.Prover.SetupHeadersForUpdate(ctx, counterparty, latestFinalizedHeader)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return headers, err
}

func (b *proverTracerBridge) CheckRefreshRequired(ctx context.Context, counterparty ChainInfoICS02Querier) (bool, error) {
	ctx, span := tracer.Start(ctx, "Prover.CheckRefreshRequired",
		trace.WithAttributes(attribute.String("counterparty_chain_id", counterparty.ChainID())),
		withPackage(b.Prover),
	)
	defer span.End()

	required, err := b.Prover.CheckRefreshRequired(ctx, counterparty)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return required, err
}

func (b *proverTracerBridge) ProveState(ctx QueryContext, path string, value []byte) ([]byte, clienttypes.Height, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Prover.ProveState",
		trace.WithAttributes(AttributeKeyPath.String(path)),
		withPackage(b.Prover),
	)
	defer span.End()

	proof, height, err := b.Prover.ProveState(ctx, path, value)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return proof, height, err
}

func (b *proverTracerBridge) ProveHostConsensusState(ctx QueryContext, height ibcexported.Height, consensusState ibcexported.ConsensusState) ([]byte, error) {
	ctx, span := StartTraceWithQueryContext(tracer, ctx, "Prover.ProveHostConsensusState",
		withPackage(b.Prover),
	)
	defer span.End()

	proof, err := b.Prover.ProveHostConsensusState(ctx, height, consensusState)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return proof, err
}

// StartTraceWithQueryContext creates a span and a QueryContext containing the newly-created span.
func StartTraceWithQueryContext(tracer trace.Tracer, ctx QueryContext, spanName string, opts ...trace.SpanStartOption) (QueryContext, trace.Span) {
	opts = append(opts, trace.WithAttributes(AttributeGroup("query",
		// Convert revision_number and revision_height to string because the attribute package does not support uint64
		AttributeKeyRevisionNumber.String(fmt.Sprint(ctx.Height().GetRevisionNumber())),
		AttributeKeyRevisionHeight.String(fmt.Sprint(ctx.Height().GetRevisionHeight())),
	)...))
	spanCtx, span := tracer.Start(ctx.Context(), spanName, opts...)
	ctx = NewQueryContext(spanCtx, ctx.Height())
	return ctx, span
}

// withPackage adds the package name of the function/method `v`
func withPackage(v any) trace.SpanStartOption {
	return trace.WithAttributes(AttributeKeyPackage.String(getPackageName(v)))
}

func getPackageName(v any) string {
	if v == nil {
		return ""
	}

	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.PkgPath()
}
