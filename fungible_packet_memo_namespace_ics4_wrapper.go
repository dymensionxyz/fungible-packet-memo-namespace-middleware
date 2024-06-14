package namespacer

import (
	"errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
)

type Transformer interface {
	Wrap(string) string
	Unwrap(string) string
}

// FungiblePacketMemoICS4Wrapper will take any fungible packet
// and namespace the memo aco
type FungiblePacketMemoICS4Wrapper struct {
	porttypes.ICS4Wrapper
	transformer Transformer
	codec       *codec.ProtoCodec
}

func NewFungiblePacketMemoNamespacingICS4Wrapper(
	next porttypes.ICS4Wrapper,
	transformer Transformer,
) *FungiblePacketMemoICS4Wrapper {
	return &FungiblePacketMemoICS4Wrapper{
		ICS4Wrapper: next,
		transformer: transformer,
	}
}

// SendPacket wraps IBC ChannelKeeper's SendPacket function. It will wrap
// the memo of any fungible token packet with the transformer and pass on the packet.
func (m *FungiblePacketMemoICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	srcPort string, srcChan string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	packet := new(transfertypes.FungibleTokenPacketData)
	if err = m.codec.UnmarshalJSON(data, packet); err != nil {
		return 0, errorsmod.Wrap(errors.Join(sdkerrors.ErrJSONUnmarshal, err), "to fungible token packet data")
	}
	return m.ICS4Wrapper.SendPacket(ctx, chanCap, srcPort, srcChan, timeoutHeight, timeoutTimestamp, data)
}
