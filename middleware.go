package memo

import (
	"errors"

	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v6/modules/core/exported"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
)

type Middlewares struct {
	Wrapper   *Wrapper
	Unwrapper *Unwrapper
}

func NewMiddlewares(
	ics4Wrapper porttypes.ICS4Wrapper,
	ibcModule porttypes.IBCModule,
	codec *codec.ProtoCodec,
	wrap, unwrap Transformer,
) Middlewares {
	return Middlewares{
		Wrapper:   NewWrapper(ics4Wrapper, codec, wrap),
		Unwrapper: NewUnwrapper(ibcModule, codec, unwrap),
	}
}

func NewDefaultMiddlewares(
	ics4Wrapper porttypes.ICS4Wrapper,
	ibcModule porttypes.IBCModule,
	codec *codec.ProtoCodec,
	namespace string,
) Middlewares {
	return NewMiddlewares(
		ics4Wrapper, ibcModule, codec, toJson(namespace), fromJson(namespace),
	)
}

type Transformer func(string) string

// Wrapper will take any fungible packet
// and namespace the memo according to the transformer
type Wrapper struct {
	porttypes.ICS4Wrapper
	wrap  Transformer
	codec *codec.ProtoCodec
}

func NewWrapper(
	next porttypes.ICS4Wrapper,
	codec *codec.ProtoCodec,
	wrap Transformer,
) *Wrapper {
	return &Wrapper{
		ICS4Wrapper: next,
		wrap:        wrap,
		codec:       codec,
	}
}

// SendPacket wraps IBC ChannelKeeper's SendPacket function. It will wrap
// the memo of any fungible token packet with the wrap and pass on the packet.
func (w *Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	srcPort string, srcChan string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	f := new(transfertypes.FungibleTokenPacketData)
	if err = w.codec.UnmarshalJSON(data, f); err != nil {
		return 0, errorsmod.Wrap(errors.Join(sdkerrors.ErrJSONUnmarshal, err), "packet data to fungible token data")
	}
	f.Memo = w.wrap(f.GetMemo())
	data, err = w.codec.MarshalJSON(f)
	if err != nil {
		return 0, errorsmod.Wrap(errors.Join(sdkerrors.ErrJSONMarshal, err), "fungible token data to packet data")
	}
	return w.ICS4Wrapper.SendPacket(ctx, chanCap, srcPort, srcChan, timeoutHeight, timeoutTimestamp, data)
}

// Unwrapper will take any fungible packet
// with a wrapped memo and unwrap it
type Unwrapper struct {
	porttypes.IBCModule
	unwrap Transformer
	codec  *codec.ProtoCodec
}

func NewUnwrapper(
	next porttypes.IBCModule,
	codec *codec.ProtoCodec,
	unwrap Transformer,
) *Unwrapper {
	return &Unwrapper{
		IBCModule: next,
		unwrap:    unwrap,
		codec:     codec,
	}
}

func (w Unwrapper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	f := new(transfertypes.FungibleTokenPacketData)
	if err := w.codec.UnmarshalJSON(packet.GetData(), f); err != nil {
		err = errorsmod.Wrap(errors.Join(sdkerrors.ErrJSONUnmarshal, err), "packet data to fungible token data")
		return channeltypes.NewErrorAcknowledgement(err)
	}
	f.Memo = w.unwrap(f.GetMemo())
	bz, err := w.codec.MarshalJSON(f)
	if err != nil {
		err = errorsmod.Wrap(errors.Join(sdkerrors.ErrJSONMarshal, err), "fungible token data to packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}
	packet.Data = bz
	return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
}
