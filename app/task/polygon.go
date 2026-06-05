package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func polygonInit() {
	ctx := context.Background()
	pol := evm{
		Network: conf.Polygon,
		Block: block{
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: pol.blockDispatch})
	Register(Task{Callback: pol.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: pol.tradeConfirmHandle, Duration: time.Second * 5})
}
