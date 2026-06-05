package admin

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/chunkburst/PreUSDT/app/handler/base"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/chunkburst/PreUSDT/app/notifier"
	"github.com/chunkburst/PreUSDT/app/utils"
	"github.com/gin-gonic/gin"
)

type Conf struct {
}

type confReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

type confGetsReq struct {
	Keys []string `json:"keys" binding:"required"`
}

type confSetsReq []confReq

type notifierConf struct {
	Channel string          `json:"channel" binding:"required"`
	Params  json.RawMessage `json:"params" binding:"required"`
}

func validateConfValue(key model.ConfKey, value string) error {
	switch key {
	case model.PaymentStaticPath:
		if value != "" && !utils.IsExist(value) {
			return fmt.Errorf("静态资源路径不存在，请确认后重新配置：%s", value)
		}
	case model.ApiAuthToken:
		return fmt.Errorf("安全考虑，不允许自定义修改 API 对接令牌")
	case model.HomepageTemplateHTML:
		if value != "" && !utils.IsSafeHomepageHTML(value) {
			return fmt.Errorf("首页 HTML 模板包含不安全内容，请检查后重试")
		}
	case model.HomepageTemplateCSS:
		if value != "" && !utils.IsSafeHomepageCSS(value) {
			return fmt.Errorf("首页 CSS 模板包含不安全内容，请检查后重试")
		}
	}

	return nil
}

func sanitizeNotifierParams(channel string, params json.RawMessage) (string, error) {
	if channel == notifier.ChannelNone {
		return "{}", nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(params, &raw); err != nil {
		return "", fmt.Errorf("通知参数格式错误")
	}

	clean := make(map[string]string)
	for key, value := range raw {
		text := strings.TrimSpace(fmt.Sprintf("%v", value))
		switch channel {
		case notifier.ChannelTelegram:
			switch key {
			case "bot_token", "chat_id", "topic_id", "success_template", "notify_fail_template", "non_order_transfer_template", "tron_resource_template", "welcome_template", "test_template":
				if strings.HasSuffix(key, "_template") && text != "" && !utils.IsSafeTelegramTemplate(text) {
					return "", fmt.Errorf("Telegram 模板 %s 包含不安全内容，请检查后重试", key)
				}
				clean[key] = text
			}
		default:
			clean[key] = text
		}
	}

	if channel == notifier.ChannelTelegram {
		if clean["bot_token"] == "" {
			return "", fmt.Errorf("Bot Token 不能为空")
		}
		if clean["chat_id"] == "" {
			return "", fmt.Errorf("Chat ID 不能为空")
		}
	}

	encoded, err := json.Marshal(clean)
	if err != nil {
		return "", fmt.Errorf("通知参数格式错误")
	}

	return string(encoded), nil
}

func (Conf) Set(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	key := model.ConfKey(strings.TrimSpace(req.Key))
	value := strings.TrimSpace(req.Value)
	if err := validateConfValue(key, value); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.SetK(key, value)

	defer model.RefreshC()

	base.Ok(ctx, "配置成功")
}

func (Conf) Get(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	base.Ok(ctx, gin.H{"key": req.Key, "value": model.GetK(model.ConfKey(req.Key))})
}

func (Conf) Del(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.Db.Where("k = ?", req.Key).Delete(&model.Conf{})

	base.Ok(ctx, "删除成功")
}

func (Conf) Gets(ctx *gin.Context) {
	var req confGetsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var items = make([]model.Conf, 0)
	model.Db.Where("k IN ?", req.Keys).Find(&items)

	var data = gin.H{}
	for _, item := range items {
		data[string(item.K)] = item.V
	}

	base.Ok(ctx, data)
}

func (Conf) Sets(ctx *gin.Context) {
	var req confSetsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	keys := make([]string, 0)
	data := make([]model.Conf, 0)
	for _, item := range req {
		var k = model.ConfKey(strings.TrimSpace(item.Key))
		var v = strings.TrimSpace(item.Value)
		if err := validateConfValue(k, v); err != nil {
			base.BadRequest(ctx, err.Error())

			return
		}
		keys = append(keys, string(k))
		data = append(data, model.Conf{K: k, V: v})
	}

	model.Db.Where("k IN ?", keys).Delete(&model.Conf{})
	model.Db.Create(&data)

	defer model.RefreshC()

	base.Ok(ctx, "配置成功")
}

func (Conf) Rpc(ctx *gin.Context) {
	var keys = []model.ConfKey{
		model.RpcEndpointPlasma,
		model.RpcEndpointBsc,
		model.RpcEndpointSolana,
		model.RpcEndpointXlayer,
		model.RpcEndpointPolygon,
		model.RpcEndpointArbitrum,
		model.RpcEndpointEthereum,
		model.RpcEndpointBase,
		model.RpcEndpointAptos,
		model.RpcEndpointTron,
		model.RpcEndpointTronGridApiKey,
	}

	var rpc = make(map[model.ConfKey]string)
	var items = make([]model.Conf, 0)
	model.Db.Where("k IN ?", keys).Find(&items)
	for _, item := range items {
		rpc[item.K] = item.V
	}

	base.Ok(ctx, gin.H{
		"rpc":   rpc,
		"stats": conf.GetStats(),
	})
}

func (Conf) Notifier(ctx *gin.Context) {
	var req notifierConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	params, err := sanitizeNotifierParams(req.Channel, req.Params)
	if err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var keys = []string{string(model.NotifierChannel), string(model.NotifierParams)}
	model.Db.Where("k IN ?", keys).Delete(&model.Conf{})
	model.Db.Create(&[]model.Conf{
		{K: model.NotifierChannel, V: req.Channel},
		{K: model.NotifierParams, V: params},
	})

	defer model.RefreshC()

	base.Ok(ctx, "配置成功")
}

func (Conf) NotifierTest(ctx *gin.Context) {
	var req notifierConf
	if err := ctx.ShouldBindJSON(&req); err == nil {
		params, err := sanitizeNotifierParams(req.Channel, req.Params)
		if err != nil {
			base.BadRequest(ctx, err.Error())

			return
		}

		n, err := notifier.NewNotifier(req.Channel, params)
		if err != nil {
			base.BadRequest(ctx, "发送测试失败："+err.Error())

			return
		}
		if err = n.Test(); err != nil {
			base.BadRequest(ctx, "发送测试失败："+err.Error())

			return
		}

		base.Ok(ctx, "发送测试成功")

		return
	}

	err := notifier.Test()
	if err != nil {
		base.BadRequest(ctx, "发送测试失败："+err.Error())

		return
	}

	base.Ok(ctx, "发送测试成功")
}

func (Conf) ResetApiAuthToken(ctx *gin.Context) {
	model.SetK(model.ApiAuthToken, strings.ToUpper(utils.Md5String(utils.StrSha256(time.Now().String()))))

	base.Ok(ctx, "重置成功")
}
