package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func baseInit() {
	ctx := context.Background()
	base := evm{
		Network: conf.Base,
		Block: block{
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: base.blockDispatch})
	Register(Task{Callback: base.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: base.tradeConfirmHandle, Duration: time.Second * 5})
}
