package epusdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/chunkburst/PreUSDT/app/handler/epay"
	"github.com/chunkburst/PreUSDT/app/log"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Epusdt struct{}

type createReq struct {
	OrderID     string     `json:"order_id" binding:"required"`
	NotifyURL   string     `json:"notify_url" binding:"required"`
	RedirectURL string     `json:"redirect_url" binding:"required"`
	Signature   string     `json:"signature" binding:"required"`
	Amount      float64    `json:"amount"`
	Name        string     `json:"name"`
	Fiat        model.Fiat `json:"fiat"`
	TradeType   string     `json:"trade_type"`
	Address     string     `json:"address"`
	Timeout     int64      `json:"timeout"`
	Rate        string     `json:"rate"`
}

type createOrderReq struct {
	OrderID     string     `json:"order_id" binding:"required"`
	NotifyURL   string     `json:"notify_url" binding:"required"`
	RedirectURL string     `json:"redirect_url" binding:"required"`
	Signature   string     `json:"signature" binding:"required"`
	Amount      float64    `json:"amount"`
	Name        string     `json:"name"`
	Fiat        model.Fiat `json:"fiat"`
	Currencies  string     `json:"currencies"`
	TradeType   string     `json:"trade_type"`
	Currency    string     `json:"currency"`
	Network     string     `json:"network"`
	Rate        string     `json:"rate"`
	Timeout     int64      `json:"timeout"`
}

type updateOrderReq struct {
	TradeID  string `json:"trade_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Network  string `json:"network" binding:"required"`
}

type cancelReq struct {
	TradeID   string `json:"trade_id" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type methodsReq struct {
	TradeID  string `json:"trade_id" binding:"required"`
	Currency string `json:"currency"`
}

type PaymentItem struct {
	Amount          string `json:"amount"`
	ActualAmount    string `json:"actual_amount"`
	Fiat            string `json:"fiat"`
	ExchangeRate    string `json:"exchange_rate"`
	Currency        string `json:"currency"`
	Network         string `json:"network"`
	TokenNetName    string `json:"token_net_name"`
	TokenCustomName string `json:"token_custom_name"`
	IsPopular       bool   `json:"is_popular"`
}

func maxSwitchCount() int {
	count := cast.ToInt(model.GetC(model.PaymentMaxSwitchCount))
	if count <= 0 {
		count = cast.ToInt("3")
	}

	return count
}

func expirationSeconds(expiredAt time.Time) uint64 {
	seconds := int64(expiredAt.Sub(time.Now()).Seconds())
	if seconds <= 0 {
		return 0
	}

	return uint64(seconds)
}

func currencyAllowedByLimit(limit string, crypto model.Crypto) bool {
	if limit == "" {
		return true
	}

	whitelist := make(map[string]bool)
	blacklist := make(map[string]bool)
	for _, item := range strings.Split(limit, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.HasPrefix(item, "-") {
			blacklist[strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(item, "-")))] = true
			continue
		}
		whitelist[strings.ToUpper(item)] = true
	}

	currency := strings.ToUpper(string(crypto))
	if blacklist[currency] {
		return false
	}
	if len(whitelist) > 0 && !whitelist[currency] {
		return false
	}

	return true
}

func orderPaymentPayload(order model.Order, host string) gin.H {
	conf, _ := model.GetTradeConfig(order.TradeType)
	maxSwitch := maxSwitchCount()
	remaining := maxSwitch - order.SwitchCount
	if remaining < 0 {
		remaining = 0
	}

	return gin.H{
		"fiat":                   order.Fiat,
		"trade_type":             order.TradeType,
		"trade_id":               order.TradeId,
		"order_id":               order.OrderId,
		"status":                 order.Status,
		"amount":                 order.Money,
		"actual_amount":          order.Amount,
		"address":                order.Address,
		"token":                  order.Address,
		"currency":               order.Crypto,
		"network":                conf.Network,
		"token_net_name":         conf.NetworkName,
		"exchange_rate":          order.Rate,
		"switch_count":           order.SwitchCount,
		"max_switch_count":       maxSwitch,
		"remaining_switch_count": remaining,
		"expiration_time":        expirationSeconds(order.ExpiredAt),
		"payment_url":            model.CheckoutCounter(host, order.TradeId),
	}
}

func paymentCurrencyRank(currency string) int {
	switch model.Crypto(currency) {
	case model.USDT:
		return 0
	case model.USDC:
		return 1
	case model.BNB:
		return 2
	case model.ETH:
		return 3
	case model.TRX:
		return 4
	default:
		return 100
	}
}

func paymentNetworkRank(network string) int {
	switch network {
	case "tron":
		return 0
	case "bsc":
		return 1
	case "ethereum":
		return 2
	case "base":
		return 3
	case "polygon":
		return 4
	case "arbitrum":
		return 5
	case "xlayer":
		return 6
	case "plasma":
		return 7
	case "solana":
		return 8
	case "aptos":
		return 9
	default:
		return 100
	}
}

func (Epusdt) SignVerify(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("json 数据读取错误 %s", err.Error())))
		ctx.Abort()

		return
	}

	m := make(map[string]any)
	if err = json.Unmarshal(rawData, &m); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("json 数据解析错误 %s", err.Error())))
		ctx.Abort()

		return
	}

	sign, ok := m["signature"]
	if !ok {
		ctx.JSON(200, respFailJson("签名丢失"))
		ctx.Abort()

		return
	}

	if utils.EpusdtSign(m, model.AuthToken()) != sign {
		ctx.JSON(200, respFailJson("签名错误"))
		ctx.Abort()

		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawData)) // 回写数据
	ctx.Next()
}

func (Epusdt) Notify(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		ctx.String(200, "fail")
		return
	}

	m := make(map[string]any)
	if err = json.Unmarshal(rawData, &m); err != nil {
		ctx.String(200, "fail")
		return
	}

	sign, ok := m["signature"]
	if !ok {
		ctx.String(200, "fail")
		return
	}

	if utils.EpusdtSign(m, model.AuthToken()) != sign {
		ctx.String(200, "fail")
		return
	}

	ctx.String(200, "ok")
}

// CreateOrder Order creation API
func (Epusdt) CreateOrder(ctx *gin.Context) {
	var req createOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("CreateOrder: request error: %s", err.Error())))

		return
	}

	if !utils.IsAllowedCallbackURL(req.NotifyURL) {
		ctx.JSON(200, respFailJson("notify_url 地址不合法"))

		return
	}
	if !utils.IsAllowedCallbackURL(req.RedirectURL) {
		ctx.JSON(200, respFailJson("redirect_url 地址不合法"))

		return
	}

	// 解析请求地址
	host := "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	if req.Fiat == "" {
		req.Fiat = model.CNY
	}

	params := model.OrderParams{
		Money:         decimal.NewFromFloat(req.Amount),
		ApiType:       model.OrderApiTypeEpusdt,
		OrderId:       req.OrderID,
		RedirectUrl:   req.RedirectURL,
		NotifyUrl:     req.NotifyURL,
		Name:          req.Name,
		Timeout:       req.Timeout,
		Rate:          req.Rate,
		Fiat:          req.Fiat,
		CurrencyLimit: req.Currencies,
		AddressLocked: req.Amount == 0,
	}

	tradeType := model.TradeType(strings.TrimSpace(req.TradeType))
	currency := strings.TrimSpace(req.Currency)
	network := strings.TrimSpace(req.Network)
	if tradeType == "" && (currency != "" || network != "") {
		if currency == "" || network == "" {
			ctx.JSON(200, respFailJson("CreateOrder: currency 和 network 必须同时提供"))
			return
		}
		parsed, err := model.GetTradeTypeByCurrencyAndNetwork(strings.ToUpper(currency), network)
		if err != nil {
			ctx.JSON(200, respFailJson(fmt.Sprintf("CreateOrder: unsupported payment method: %s - %s", req.Currency, req.Network)))
			return
		}
		tradeType = parsed
	}

	var order model.Order
	var err error
	if tradeType != "" {
		conf, ok := model.GetTradeConfig(tradeType)
		if !ok {
			ctx.JSON(200, respFailJson(fmt.Sprintf("CreateOrder: unsupported trade type: %s", tradeType)))
			return
		}
		if !currencyAllowedByLimit(req.Currencies, conf.Crypto) {
			ctx.JSON(200, respFailJson("CreateOrder: 当前订单不允许使用该币种"))
			return
		}
		params.TradeType = tradeType
		order, err = model.StartBuildOrder(params)
	} else {
		order, err = model.BuildPendingOrder(params)
	}
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("CreateOrder: order create failed: %s", err.Error())))
		return
	}

	log.Info(fmt.Sprintf("订单创建成功 商户订单：%s", req.OrderID))

	// 返回响应数据
	ctx.JSON(200, respSuccJson(gin.H{
		"fiat":            order.Fiat,
		"trade_id":        order.TradeId,
		"order_id":        order.OrderId,
		"name":            order.Name,
		"status":          order.Status,
		"amount":          order.Money,
		"expiration_time": expirationSeconds(order.ExpiredAt),
		"payment_url":     model.CheckoutCashier(host, order.TradeId),
		"network":         GetPaymentItem("", order),
	}))
}

// UpdateOrder 更新订单支付方式，返回当前页可直接刷新的付款信息。
func (Epusdt) UpdateOrder(ctx *gin.Context) {
	var req updateOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("UpdateOrder: request error: %s", err.Error())))
		return
	}

	host := utils.GetRequestHost(ctx.Request)

	order, ok := model.GetTradeOrder(req.TradeID)
	if !ok {
		ctx.JSON(200, respFailJson("order not found"))
		return
	}

	if order.Status != model.OrderStatusWaiting {
		ctx.JSON(200, respFailJson("当前订单状态不允许切换支付方式"))
		return
	}
	if !time.Now().Before(order.ExpiredAt) {
		ctx.JSON(200, respFailJson("订单已过期"))
		return
	}

	tradeType, err := model.GetTradeTypeByCurrencyAndNetwork(req.Currency, req.Network)
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("UpdateOrder: unsupported payment method: %s - %s", req.Currency, req.Network)))
		return
	}
	conf, _ := model.GetTradeConfig(tradeType)
	if !currencyAllowedByLimit(order.CurrencyLimit, conf.Crypto) {
		ctx.JSON(200, respFailJson("当前订单不允许使用该币种"))
		return
	}

	if tradeType == order.TradeType {
		ctx.JSON(200, respSuccJson(orderPaymentPayload(order, host)))
		return
	}

	initialSelection := order.Address == "" || order.Amount == "0"
	maxSwitch := maxSwitchCount()
	if !initialSelection && order.SwitchCount >= maxSwitch {
		ctx.JSON(200, respFailJson("订单支付方式切换次数已用完"))
		return
	}

	money, _ := decimal.NewFromString(order.Money)
	params := model.OrderParams{
		Money:         money,
		OrderId:       order.OrderId,
		TradeType:     tradeType,
		RedirectUrl:   order.ReturnUrl,
		NotifyUrl:     order.NotifyUrl,
		Name:          order.Name,
		Timeout:       int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
		Fiat:          order.Fiat,
		AddressLocked: money.IsZero(),
	}

	newOrder, err := model.SwitchOrderPayment(order, params, initialSelection, maxSwitch)
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("update order failed: %s", err.Error())))
		return
	}

	ctx.JSON(200, respSuccJson(orderPaymentPayload(newOrder, host)))
}

func (Epusdt) CreateTransaction(ctx *gin.Context) {
	var req createReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("请求参数错误：%s", err.Error())))

		return
	}

	if !utils.IsAllowedCallbackURL(req.NotifyURL) {
		ctx.JSON(200, respFailJson("notify_url 地址不合法"))

		return
	}
	if !utils.IsAllowedCallbackURL(req.RedirectURL) {
		ctx.JSON(200, respFailJson("redirect_url 地址不合法"))

		return
	}

	if req.Fiat == "" {
		req.Fiat = model.CNY
	}
	if req.TradeType == "" {
		req.TradeType = string(model.UsdtTrc20)
	}

	order, err := model.StartBuildOrder(model.OrderParams{
		Money:         decimal.NewFromFloat(req.Amount),
		ApiType:       model.OrderApiTypeEpusdt,
		Address:       req.Address,
		AddressLocked: req.Amount == 0,
		OrderId:       req.OrderID,
		TradeType:     model.TradeType(req.TradeType),
		RedirectUrl:   req.RedirectURL,
		NotifyUrl:     req.NotifyURL,
		Name:          req.Name,
		Timeout:       req.Timeout,
		Rate:          req.Rate,
		Fiat:          req.Fiat,
	})
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("订单创建失败：%s", err.Error())))

		return
	}

	log.Info(fmt.Sprintf("订单创建成功 商户订单：%s", req.OrderID))

	// 返回响应数据
	ctx.JSON(200, respSuccJson(gin.H{
		"fiat":            order.Fiat,
		"trade_type":      order.TradeType,
		"trade_id":        order.TradeId,
		"order_id":        order.OrderId,
		"status":          order.Status,
		"amount":          order.Money,
		"actual_amount":   order.Amount,
		"token":           order.Address,
		"expiration_time": expirationSeconds(order.ExpiredAt),
		"payment_url":     model.CheckoutCounter(utils.GetRequestHost(ctx.Request), order.TradeId),
	}))
}

func (Epusdt) CancelTransaction(ctx *gin.Context) {
	var req cancelReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("请求参数错误：%s", err.Error())))

		return
	}

	order, ok := model.GetTradeOrder(req.TradeID)
	if !ok {
		ctx.JSON(200, respFailJson("订单不存在"))

		return
	}

	if order.Status != model.OrderStatusWaiting {
		ctx.JSON(200, respFailJson(fmt.Sprintf("当前订单(%s)状态不允许取消", req.TradeID)))

		return
	}

	if err := order.SetCanceled(); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("订单取消失败：%s", err.Error())))

		return
	}

	ctx.JSON(200, respSuccJson(gin.H{"trade_id": req.TradeID}))
}

func (Epusdt) CheckoutCounter(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.String(200, "订单不存在")

		return
	}

	uri, err := url.ParseRequestURI(order.ReturnUrl)
	if err != nil {
		ctx.String(200, "同步地址错误")
		log.Error("同步地址解析错误", err.Error())

		return
	}

	ctx.HTML(200, string(order.TradeType+".html"), gin.H{
		"http_host":  uri.Host,
		"amount":     order.Amount,
		"address":    order.Address,
		"expire":     int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
		"return_url": order.ReturnUrl,
		"usdt_rate":  order.Rate,
		"trade_id":   tradeId,
		"order_id":   order.OrderId,
		"trade_type": order.TradeType,
		"money":      order.Money,
		"fiat":       order.Fiat,
	})
}

func (Epusdt) CheckoutCashier(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.String(200, "order not found")

		return
	}

	uri, err := url.ParseRequestURI(order.ReturnUrl)
	if err != nil {
		ctx.String(200, "sync address error")
		log.Error("CheckoutCashier: sync address error: ", err.Error())

		return
	}

	payment := orderPaymentPayload(order, utils.GetRequestHost(ctx.Request))

	ctx.HTML(200, "cashier.html", gin.H{
		"http_host":              uri.Host,
		"amount":                 order.Amount,
		"expire":                 int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
		"return_url":             order.ReturnUrl,
		"trade_id":               tradeId,
		"order_id":               order.OrderId,
		"name":                   order.Name,
		"money":                  order.Money,
		"fiat":                   order.Fiat,
		"trade_type":             order.TradeType,
		"address":                order.Address,
		"crypto":                 order.Crypto,
		"rate":                   order.Rate,
		"switch_count":           order.SwitchCount,
		"max_switch_count":       payment["max_switch_count"],
		"remaining_switch_count": payment["remaining_switch_count"],
		"current_payment":        payment,
		"network":                GetPaymentItem("", order),
	})
}

func (Epusdt) CheckStatus(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.JSON(200, respFailJson("订单不存在"))

		return
	}

	var returnUrl string
	if order.Status == model.OrderStatusSuccess {
		returnUrl = order.ReturnUrl
		if order.ApiType == model.OrderApiTypeEpay {
			// 易支付兼容
			returnUrl = fmt.Sprintf("%s?%s", returnUrl, epay.BuildNotifyParams(order))
		}
	}

	// 返回响应数据
	ctx.JSON(200, gin.H{
		"trade_id":   tradeId,
		"trade_hash": order.RefHash,
		"status":     order.Status,
		"return_url": returnUrl,
	})
}

func (Epusdt) GetPaymentMethods(ctx *gin.Context) {
	var req methodsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(500, respFailJson(fmt.Sprintf("request error: %s", err.Error())))
		return
	}

	order, ok := model.GetTradeOrder(req.TradeID)
	if !ok {
		ctx.JSON(400, respFailJson("order not found"))
		return
	}

	if order.Status == model.OrderStatusExpired {
		ctx.JSON(400, respFailJson("order expired"))
		return
	}

	ctx.JSON(200, respSuccJson(gin.H{
		"methods": GetPaymentItem(model.Crypto(req.Currency), order),
	}))
}

func GetPaymentItem(crypto model.Crypto, order model.Order) []PaymentItem {
	fiat := order.Fiat

	var methods = make([]PaymentItem, 0)
	allTrades := model.GetAllTradeConfig()

	// 解析限定币种
	var whitelist = make(map[string]bool)
	var blacklist = make(map[string]bool)
	if order.CurrencyLimit != "" {
		for _, c := range strings.Split(order.CurrencyLimit, ",") {
			c = strings.TrimSpace(c)
			if strings.HasPrefix(c, "-") {
				blacklist[strings.ToUpper(strings.TrimPrefix(c, "-"))] = true
			} else {
				whitelist[strings.ToUpper(c)] = true
			}
		}
	}

	for tradeTypeStr, conf := range allTrades {
		// 如果指定了货币，则进行过滤
		if crypto != "" && conf.Crypto != crypto {
			continue
		}

		// Check blacklist
		if len(blacklist) > 0 && blacklist[string(conf.Crypto)] {
			continue
		}

		// Check whitelist
		if len(whitelist) > 0 && !whitelist[string(conf.Crypto)] {
			continue
		}

		// 检查是否有可用钱包
		count := len(model.GetAvailableAddress(model.TradeType(tradeTypeStr)))
		if count == 0 {
			continue
		}

		// 获取汇率配置的浮动语法
		syntax := model.GetK(model.ConfKey(fmt.Sprintf("rate_float_%s_%s", conf.Crypto, fiat)))

		// 获取汇率
		rate, err := model.GetOrderRate(conf.Crypto, fiat, syntax)
		if err != nil {
			log.Error(fmt.Sprintf("GetPaymentMethods: get order rate error: %s", err.Error()))
			continue
		}

		// 计算实际支付金额 (加密货币)
		// Money 是法币金额
		moneyDecimal, _ := decimal.NewFromString(order.Money)

		// 计算精度
		atom, precision := model.GetAtomicity(model.TradeType(tradeTypeStr))
		actualAmount := moneyDecimal.DivRound(rate, precision)
		if actualAmount.LessThan(atom) {
			actualAmount = atom
		}

		methods = append(methods, PaymentItem{
			Amount:          order.Money,
			ActualAmount:    actualAmount.String(),
			Fiat:            string(fiat),
			ExchangeRate:    rate.String(),
			Currency:        string(conf.Crypto),
			Network:         string(conf.Network),
			TokenNetName:    conf.NetworkName,
			TokenCustomName: "",    // 暂为空
			IsPopular:       false, // 暂为 false
		})
	}

	sort.Slice(methods, func(i, j int) bool {
		leftCurrencyRank := paymentCurrencyRank(methods[i].Currency)
		rightCurrencyRank := paymentCurrencyRank(methods[j].Currency)
		if leftCurrencyRank != rightCurrencyRank {
			return leftCurrencyRank < rightCurrencyRank
		}
		if methods[i].Currency != methods[j].Currency {
			return methods[i].Currency < methods[j].Currency
		}

		leftNetworkRank := paymentNetworkRank(methods[i].Network)
		rightNetworkRank := paymentNetworkRank(methods[j].Network)
		if leftNetworkRank != rightNetworkRank {
			return leftNetworkRank < rightNetworkRank
		}

		return methods[i].Network < methods[j].Network
	})

	return methods
}

func respFailJson(message string) gin.H {

	return gin.H{"status_code": 400, "message": message}
}

func respSuccJson(data interface{}) gin.H {

	return gin.H{"status_code": 200, "message": "success", "data": data, "request_id": ""}
}
