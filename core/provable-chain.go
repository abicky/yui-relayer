package core

import (
	"context"
	"reflect"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ProvableChain represents a chain that is supported by the relayer.
type ProvableChain struct {
	Chain
	Prover
}

// NewProvableChain returns a new ProvableChain instance
func NewProvableChain(chain Chain, prover Prover) *ProvableChain {
	return &ProvableChain{Chain: chain, Prover: prover}
}

// As is a method similar to errors.As and finds the first struct value in the specified
// field that matches target.
//
// In the following example, you can set a struct value in the Chain field to the chain variable:
//
//	var chain module.Chain
//	if ok := provableChain.As("Chain", &chain); !ok {
//		return errors.New("Chain is not a module.Chain")
//	}
func (pc *ProvableChain) As(fieldName string, v any) bool {
	targetType := reflect.TypeOf(v)

	rv := reflect.ValueOf(pc)
	for {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		field := rv.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}

		fieldValue := field.Interface()
		rv = reflect.ValueOf(fieldValue)
		if reflect.TypeOf(fieldValue).AssignableTo(targetType) {
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			reflect.ValueOf(v).Elem().Set(rv)
			return true
		}
	}
}

func (pc *ProvableChain) Init(homePath string, timeout time.Duration, codec codec.ProtoCodecMarshaler, debug bool) error {
	if err := pc.Chain.Init(homePath, timeout, codec, debug); err != nil {
		return err
	}
	if err := pc.Prover.Init(homePath, timeout, codec, debug); err != nil {
		return err
	}
	return nil
}

func (pc *ProvableChain) SetRelayInfo(path *PathEnd, counterparty *ProvableChain, counterpartyPath *PathEnd) error {
	if err := pc.Chain.SetRelayInfo(path, counterparty, counterpartyPath); err != nil {
		return err
	}
	if err := pc.Prover.SetRelayInfo(path, counterparty, counterpartyPath); err != nil {
		return err
	}
	return nil
}

func (pc *ProvableChain) SetupForRelay(ctx context.Context) error {
	if err := pc.Chain.SetupForRelay(ctx); err != nil {
		return err
	}
	if err := pc.Prover.SetupForRelay(ctx); err != nil {
		return err
	}
	return nil
}
