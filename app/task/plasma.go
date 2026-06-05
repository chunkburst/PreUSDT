package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func plasmaInit() {
	ctx := context.Background()
	xpl := evm{
		Network: conf.Plasma,
		Block: block{
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: xpl.blockDispatch})
	Register(Task{Callback: xpl.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: xpl.tradeConfirmHandle, Duration: time.Second * 5})
}
