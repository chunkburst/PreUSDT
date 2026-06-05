package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func arbitrumInit() {
	ctx := context.Background()
	arb := evm{
		Network: conf.Arbitrum,
		Block: block{
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: arb.blockDispatch})
	Register(Task{Callback: arb.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: arb.tradeConfirmHandle, Duration: time.Second * 5})
}
