# memo

The memo package provides middleware functionality for handling memos in fungible token packets in the Inter-Blockchain Communication (IBC) protocol. The package allows for the wrapping and unwrapping of memos within IBC packets using customizable transformation functions.

## Usage

The package provides two main middleware structures: Wrapper and Unwrapper. These can be initialized using the NewMiddlewares or NewDefaultMiddlewares functions.

```
import (
    "github.com/your-repo/memo"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
)

codec := codec.NewProtoCodec()
ics4Wrapper := // your ICS4Wrapper instance
ibcModule := // your IBCModule instance
namespace := "yourNamespace" # prefer short

middlewares := memo.NewDefaultMiddlewares(ics4Wrapper, ibcModule, codec, namespace)
```