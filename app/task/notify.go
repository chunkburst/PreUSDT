package task

import (
	"context"
	"time"

	"github.com/chunkburst/PreUSDT/app/log"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/notifier"
	"github.com/chunkburst/PreUSDT/app/task/notify"
	"github.com/chunkburst/PreUSDT/app/utils"
)

func init() {
	Register(Task{Duration: time.Second * 3, Callback: notifyRetry})
	Register(Task{Duration: time.Second * 30, Callback: notifyRoll})
}

// notifyRetry 回调失败重试
func notifyRetry(context.Context) {
	tradeOrders, err := model.GetNotifyFailedTradeOrders()
	if err != nil {
		log.Task.Error("待回调订单获取失败", err)

		return
	}

	for _, order := range tradeOrders {
		next := utils.CalcNextNotifyTime(*order.ConfirmedAt, order.NotifyNum)
		if time.Now().Unix() >= next.Unix() {
			go notify.Handle(order)
		}
	}
}

func notifyRoll(context.Context) {
	for _, o := range model.GetOrderByStatus(model.OrderStatusWaiting) {
		notify.PreUSDT(o)
	}
}

// notifyOrderSuccess 统一触发订单成功后的回调与订单通知。
func notifyOrderSuccess(order model.Order) {
	go notify.Handle(order)
	go notifier.Success(order)
}
