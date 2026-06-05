package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func bscInit() {
	ctx := context.Background()
	bsc := evm{
		Network: conf.Bsc,
		Block: block{
			ConfirmedOffset: 15,
		},
		Native: evmNative{
			Parse:     true,
			Decimal:   conf.BscBnbDecimals,
			TradeType: model.BscBnb,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: bsc.blockDispatch})
	Register(Task{Callback: bsc.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: bsc.tradeConfirmHandle, Duration: time.Second * 5})
}
