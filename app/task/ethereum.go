package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/smallnest/chanx"
)

func ethInit() {
	ctx := context.Background()
	eth := evm{
		Network: conf.Ethereum,
		Block: block{
			ConfirmedOffset: 12,
		},
		Native: evmNative{
			Parse:     true,
			TradeType: model.EthereumEth,
			Decimal:   conf.EthereumEthDecimals,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: eth.blockDispatch})
	Register(Task{Callback: eth.syncBlocksForward, Duration: time.Second * 12})
	Register(Task{Callback: eth.tradeConfirmHandle, Duration: time.Second * 5})
}
