package core

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// ProvableChain represents a chain that is supported by the relayer
type ProvableChain struct {
	Chain
	Prover
}

// NewProvableChain returns a new ProvableChain instance
func NewProvableChain(chain Chain, prover Prover) *ProvableChain {
	return &ProvableChain{Chain: chain, Prover: prover}
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

func (pc *ProvableChain) SendMsgs(ctx context.Context, msgs []sdk.Msg) ([]MsgID, error) {
	return pc.Chain.SendMsgs(ctx, msgs)
}

func (pc *ProvableChain) GetMsgResult(ctx context.Context, id MsgID) (MsgResult, error) {
	return pc.Chain.GetMsgResult(ctx, id)
}

func (pc *ProvableChain) LatestHeight(ctx context.Context) (ibcexported.Height, error) {
	return pc.Chain.LatestHeight(ctx)
}

func (pc *ProvableChain) Timestamp(ctx context.Context, height ibcexported.Height) (time.Time, error) {
	return pc.Chain.Timestamp(ctx, height)
}

func (pc *ProvableChain) QueryClientConsensusState(ctx QueryContext, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	return pc.Chain.QueryClientConsensusState(ctx, dstClientConsHeight)
}

func (pc *ProvableChain) QueryClientState(ctx QueryContext) (*clienttypes.QueryClientStateResponse, error) {
	return pc.Chain.QueryClientState(ctx)
}

func (pc *ProvableChain) QueryConnection(ctx QueryContext, connectionID string) (*conntypes.QueryConnectionResponse, error) {
	return pc.Chain.QueryConnection(ctx, connectionID)
}

func (pc *ProvableChain) QueryChannel(ctx QueryContext) (chanRes *chantypes.QueryChannelResponse, err error) {
	return pc.Chain.QueryChannel(ctx)
}

func (pc *ProvableChain) QueryUnreceivedPackets(ctx QueryContext, seqs []uint64) ([]uint64, error) {
	return pc.Chain.QueryUnreceivedPackets(ctx, seqs)
}

func (pc *ProvableChain) QueryUnfinalizedRelayPackets(ctx QueryContext, counterparty LightClientICS04Querier) (PacketInfoList, error) {
	return pc.Chain.QueryUnfinalizedRelayPackets(ctx, counterparty)
}

func (pc *ProvableChain) QueryUnreceivedAcknowledgements(ctx QueryContext, seqs []uint64) ([]uint64, error) {
	return pc.Chain.QueryUnreceivedAcknowledgements(ctx, seqs)
}

func (pc *ProvableChain) QueryUnfinalizedRelayAcknowledgements(ctx QueryContext, counterparty LightClientICS04Querier) (PacketInfoList, error) {
	return pc.Chain.QueryUnfinalizedRelayAcknowledgements(ctx, counterparty)
}

func (pc *ProvableChain) QueryChannelUpgrade(ctx QueryContext) (*chantypes.QueryUpgradeResponse, error) {
	return pc.Chain.QueryChannelUpgrade(ctx)
}

func (pc *ProvableChain) QueryChannelUpgradeError(ctx QueryContext) (*chantypes.QueryUpgradeErrorResponse, error) {
	return pc.Chain.QueryChannelUpgradeError(ctx)
}

func (pc *ProvableChain) QueryCanTransitionToFlushComplete(ctx QueryContext) (bool, error) {
	return pc.Chain.QueryCanTransitionToFlushComplete(ctx)
}

func (pc *ProvableChain) QueryBalance(ctx QueryContext, address sdk.AccAddress) (sdk.Coins, error) {
	return pc.Chain.QueryBalance(ctx, address)
}

func (pc *ProvableChain) QueryDenomTraces(ctx QueryContext, offset, limit uint64) (*transfertypes.QueryDenomTracesResponse, error) {
	return pc.Chain.QueryDenomTraces(ctx, offset, limit)
}

func (pc *ProvableChain) GetLatestFinalizedHeader(ctx context.Context) (latestFinalizedHeader Header, err error) {
	return pc.Prover.GetLatestFinalizedHeader(ctx)
}

func (pc *ProvableChain) CreateInitialLightClientState(ctx context.Context, height ibcexported.Height) (ibcexported.ClientState, ibcexported.ConsensusState, error) {
	return pc.Prover.CreateInitialLightClientState(ctx, height)
}

func (pc *ProvableChain) SetupHeadersForUpdate(ctx context.Context, counterparty FinalityAwareChain, latestFinalizedHeader Header) ([]Header, error) {
	return pc.Prover.SetupHeadersForUpdate(ctx, counterparty, latestFinalizedHeader)
}

func (pc *ProvableChain) CheckRefreshRequired(ctx context.Context, counterparty ChainInfoICS02Querier) (bool, error) {
	return pc.Prover.CheckRefreshRequired(ctx, counterparty)
}

func (pc *ProvableChain) ProveState(ctx QueryContext, path string, value []byte) (proof []byte, proofHeight clienttypes.Height, err error) {
	return pc.Prover.ProveState(ctx, path, value)
}

func (pc *ProvableChain) ProveHostConsensusState(ctx QueryContext, height ibcexported.Height, consensusState ibcexported.ConsensusState) (proof []byte, err error) {
	return pc.Prover.ProveHostConsensusState(ctx, height, consensusState)
}
