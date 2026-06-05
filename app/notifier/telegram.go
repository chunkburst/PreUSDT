package notifier

import (
	"context"
	"strings"
	"time"

	"github.com/chunkburst/PreUSDT/app"
	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/log"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/v03413/tronprotocol/core"
)

const (
	telegramSuccessTemplateKey          = "success_template"
	telegramNotifyFailTemplateKey       = "notify_fail_template"
	telegramNonOrderTransferTemplateKey = "non_order_transfer_template"
	telegramTronResourceTemplateKey     = "tron_resource_template"
	telegramWelcomeTemplateKey          = "welcome_template"
	telegramTestTemplateKey             = "test_template"
)

type Telegram struct {
	api     *bot.Bot
	token   string
	chatID  int64
	topicID int
	params  gjson.Result
}

func (t *Telegram) Initialize(params string) error {
	info := gjson.Parse(params)
	b, err := bot.New(info.Get("bot_token").String())
	if err != nil {
		return err
	}

	t.token = info.Get("bot_token").String()
	t.chatID = info.Get("chat_id").Int()
	t.topicID = cast.ToInt(info.Get("topic_id").Int())
	t.params = info
	t.api = b

	return nil
}

func (t *Telegram) template(key, fallback string) string {
	value := strings.TrimSpace(t.params.Get(key).String())
	if value == "" {
		return fallback
	}

	return value
}

func (t *Telegram) renderTemplate(key, fallback string, values map[string]string) string {
	text := t.template(key, fallback)
	for placeholder, value := range values {
		text = strings.ReplaceAll(text, "{{"+placeholder+"}}", value)
	}

	return text
}

func (t *Telegram) Success(o model.Order) {
	if o.Status != model.OrderStatusSuccess {
		return
	}

	fallback := strings.TrimSpace(`
✅ 收款成功
商户订单：{{order_id}}
请求金额：{{money}} {{fiat}}（汇率 {{rate}}）
支付数额：{{amount}} {{trade_type}}
交易哈希：{{tx_hash}}
收款地址：{{address}}
创建时间：{{created_at}}
支付时间：{{updated_at}}`)

	text := t.renderTemplate(telegramSuccessTemplateKey, fallback, map[string]string{
		"order_id":   o.OrderId,
		"money":      o.Money,
		"fiat":       string(o.Fiat),
		"rate":       o.Rate,
		"amount":     o.Amount,
		"trade_type": string(o.TradeType),
		"tx_hash":    utils.MaskHash(o.RefHash),
		"address":    utils.MaskAddress(o.Address),
		"created_at": o.CreatedAt.Format(time.DateTime),
		"updated_at": o.UpdatedAt.Format(time.DateTime),
	})

	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "📝查看交易明细", URL: o.GetTxUrl()},
			}},
		},
	})
}

func (t *Telegram) NotifyFail(o model.Order, reason string) {
	confirmedAt := ""
	nextNotifyAt := ""
	if o.ConfirmedAt != nil {
		confirmedAt = o.ConfirmedAt.Format(time.DateTime)
		nextNotifyAt = utils.CalcNextNotifyTime(*o.ConfirmedAt, o.NotifyNum+1).Format(time.DateTime)
	}

	fallback := strings.TrimSpace(`
⚠️ 回调失败
商户订单：{{order_id}}
支付数额：{{amount}}
请求金额：{{money}} {{fiat}}（汇率 {{rate}}）
交易类别：{{trade_type}}
确认时间：{{confirmed_at}}
下次回调：{{next_notify_at}}
失败原因：{{reason}}`)

	text := t.renderTemplate(telegramNotifyFailTemplateKey, fallback, map[string]string{
		"order_id":       o.OrderId,
		"amount":         o.Amount,
		"money":          o.Money,
		"fiat":           string(o.Fiat),
		"rate":           o.Rate,
		"trade_type":     strings.ToUpper(string(o.TradeType)),
		"confirmed_at":   confirmedAt,
		"next_notify_at": nextNotifyAt,
		"reason":         reason,
	})

	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "📝查看收款详情", URL: o.GetTxUrl()},
			}},
		},
	})
}

func (t *Telegram) NonOrderTransfer(trans model.TronTransfer, wa model.Wallet) {
	direction := "收入"
	if trans.RecvAddress != wa.Address {
		direction = "支出"
	}

	fallback := strings.TrimSpace(`
📨 非订单交易
监控方向：{{direction}}
交易数额：{{amount}}
交易类别：{{trade_type}}
交易时间：{{timestamp}}
接收地址：{{recv_address}}
发送地址：{{from_address}}`)

	text := t.renderTemplate(telegramNonOrderTransferTemplateKey, fallback, map[string]string{
		"direction":    direction,
		"amount":       trans.Amount.String(),
		"trade_type":   strings.ToUpper(string(trans.TradeType)),
		"timestamp":    trans.Timestamp.Format(time.DateTime),
		"recv_address": utils.MaskAddress(trans.RecvAddress),
		"from_address": utils.MaskAddress(trans.FromAddress),
	})

	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "📝查看交易明细", URL: model.GetTxUrl(trans.TradeType, trans.TxHash)},
			}},
		},
	})
}

func (t *Telegram) TronResourceChange(res model.TronResource) {
	action := "代理"
	if res.Type == core.Transaction_Contract_UnDelegateResourceContract {
		action = "回收"
	}

	fallback := strings.TrimSpace(`
🔋 资源动态
操作类型：{{action}}
质押数量：{{balance}}
交易时间：{{timestamp}}
操作地址：{{recv_address}}
资源来源：{{from_address}}`)

	text := t.renderTemplate(telegramTronResourceTemplateKey, fallback, map[string]string{
		"action":       action,
		"balance":      cast.ToString(res.Balance / 1000000),
		"timestamp":    res.Timestamp.Format(time.DateTime),
		"recv_address": utils.MaskAddress(res.RecvAddress),
		"from_address": utils.MaskAddress(res.FromAddress),
	})

	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "📝查看交易明细", URL: "https://tronscan.org/#/transaction/" + res.ID},
			}},
		},
	})
}

func (t *Telegram) Welcome() {
	fallback := strings.TrimSpace(`
👋 欢迎使用 preusdt，{{desc}}，如果您看到此消息，说明系统已启动成功！
当前版本：{{version}}
开源地址：{{github}}`)

	text := t.renderTemplate(telegramWelcomeTemplateKey, fallback, map[string]string{
		"desc":    conf.Desc,
		"version": app.Version,
		"github":  conf.Github,
	})

	t.sendMessage(&bot.SendMessageParams{
		Text: text,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "📢 关注频道", URL: "https://t.me/PreUSDTChannel"},
				{Text: "💬 社区交流", URL: "https://t.me/PreUSDTChat"},
			}},
		},
	})
}

func (t *Telegram) Test() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	text := t.renderTemplate(telegramTestTemplateKey, strings.TrimSpace(`
✅ 这是一条测试消息，Telegram 通知配置成功！
当前系统时间：{{now}}`), map[string]string{
		"now": time.Now().Format("2006-01-02 15:04:05"),
	})

	_, err := t.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          t.chatID,
		MessageThreadID: t.topicID,
		Text:            text,
	})

	return err
}

func (t *Telegram) sendMessage(p *bot.SendMessageParams) {
	p.ChatID = t.chatID
	p.MessageThreadID = t.topicID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := t.api.SendMessage(ctx, p)
	if err != nil {
		log.Warn("Bot Send Message Error:", err.Error())
	}
}
