package memo

import (
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func TestDefaultHappyPath(t *testing.T) {
	codec := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())

	next := mockMiddlewares{}
	m := NewDefaultMiddlewares(next, next, codec, "foo")

	memo := "bar"

	f := new(transfertypes.FungibleTokenPacketData)
	f.Memo = memo

	bz, _ := codec.MarshalJSON(f)

	m.Wrapper.SendPacket(sdk.Context{}, nil, "", "", clienttypes.Height{}, 0, bz)
	m.Unwrapper.OnRecvPacket(sdk.Context{}, channeltypes.Packet{Data: next.sent}, nil)
	got := next.received.GetData()

	_ = codec.UnmarshalJSON(got, f)
	require.Equal(t, memo, f.GetMemo())
}

var _ porttypes.Middleware = &mockMiddlewares{}

type mockMiddlewares struct {
	sent     []byte
	received channeltypes.Packet
}

func (m mockMiddlewares) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	m.received = packet
	return nil
}

func (m mockMiddlewares) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	m.sent = data
	return 0, nil
}

func (m mockMiddlewares) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	// TODO implement me
	panic("implement me")
}

func (m mockMiddlewares) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	// TODO implement me
	panic("implement me")
}
