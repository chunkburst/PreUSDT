package admin

import (
	"fmt"
	"strings"

	"github.com/chunkburst/PreUSDT/app/handler/base"
	"github.com/chunkburst/PreUSDT/app/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Wallet struct{}

type wAddReq struct {
	Name        string           `json:"name"`
	Remark      string           `json:"remark"`
	Address     string           `json:"address" binding:"required"`
	TradeType   string           `json:"trade_type"`
	TradeTypes  []string         `json:"trade_types"`
	Kind        model.WalletKind `json:"kind"`
	Base        model.WalletBase `json:"base"`
	OtherNotify uint8            `json:"other_notify"`
}

type wModReq struct {
	base.IDRequest
	Name        string           `json:"name"`
	Status      uint8            `json:"status"`
	Address     string           `json:"address" binding:"required"`
	Remark      string           `json:"remark"`
	TradeType   string           `json:"trade_type"`
	TradeTypes  []string         `json:"trade_types"`
	Kind        model.WalletKind `json:"kind"`
	Base        model.WalletBase `json:"base"`
	OtherNotify uint8            `json:"other_notify"`
}

type wListReq struct {
	base.ListRequest
	Name    string `json:"name"`
	Address string `json:"address"`
	Trade   string `json:"trade_type"`
}

func normalizeWalletKind(kind model.WalletKind) model.WalletKind {
	if kind == model.WalletKindMulti {
		return model.WalletKindMulti
	}

	return model.WalletKindSingle
}

func normalizeWalletBase(base model.WalletBase) model.WalletBase {
	if base == model.WalletBaseEvm {
		return model.WalletBaseEvm
	}

	return model.WalletBaseGeneric
}

func resolveWalletTradeTypes(kind model.WalletKind, baseType model.WalletBase, tradeType string, tradeTypes []string) ([]model.TradeType, error) {
	if kind == model.WalletKindMulti {
		if baseType == model.WalletBaseGeneric {
			return nil, fmt.Errorf("多链钱包必须选择基础分类")
		}
		if len(tradeTypes) == 0 {
			return nil, fmt.Errorf("多链钱包至少选择一个交易类型")
		}

		result := make([]model.TradeType, 0, len(tradeTypes))
		seen := make(map[model.TradeType]struct{}, len(tradeTypes))
		for _, item := range tradeTypes {
			t := model.TradeType(strings.TrimSpace(item))
			if t == "" {
				continue
			}
			if !model.IsSupportedTradeType(t) {
				return nil, fmt.Errorf("不支持的交易类型: %s", item)
			}
			if !model.IsTradeTypeAllowedForWalletBase(baseType, t) {
				return nil, fmt.Errorf("交易类型 %s 不属于基础分类 %s", item, baseType)
			}
			if _, ok := seen[t]; ok {
				continue
			}
			seen[t] = struct{}{}
			result = append(result, t)
		}
		if len(result) == 0 {
			return nil, fmt.Errorf("多链钱包至少选择一个有效交易类型")
		}

		return result, nil
	}

	t := model.TradeType(strings.TrimSpace(tradeType))
	if !model.IsSupportedTradeType(t) {
		return nil, fmt.Errorf("不支持的交易类型: %s", tradeType)
	}

	return []model.TradeType{t}, nil
}

func buildWalletRows(groupID string, name, remark, address string, status uint8, otherNotify uint8, kind model.WalletKind, baseType model.WalletBase, tradeTypes []model.TradeType) ([]model.Wallet, error) {
	rows := make([]model.Wallet, 0, len(tradeTypes))
	for _, tradeType := range tradeTypes {
		wallet := model.Wallet{
			Name:        strings.TrimSpace(name),
			Remark:      strings.TrimSpace(remark),
			Address:     strings.TrimSpace(address),
			TradeType:   string(tradeType),
			Status:      status,
			OtherNotify: otherNotify,
			GroupID:     groupID,
			Kind:        kind,
			Base:        baseType,
		}
		wallet.NormalizeAddress()
		wallet.FillGrouping()
		if !wallet.IsValid() {
			return nil, fmt.Errorf("钱包地址格式不合法，请检查")
		}
		rows = append(rows, wallet)
	}

	return rows, nil
}

func loadWalletGroupRows(walletID int) ([]model.Wallet, error) {
	var current model.Wallet
	if err := model.Db.Where("id = ?", walletID).First(&current).Error; err != nil {
		return nil, err
	}
	current.FillGrouping()

	var rows []model.Wallet
	if err := model.Db.Where("group_id = ? OR id = ?", current.GroupID, current.ID).Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		rows = append(rows, current)
	}

	return rows, nil
}

func deleteWalletRows(tx *gorm.DB, rows []model.Wallet) error {
	if len(rows) == 0 {
		return nil
	}

	storedGroupID := strings.TrimSpace(rows[0].GroupID)
	if storedGroupID != "" {
		return tx.Where("group_id = ?", storedGroupID).Delete(&model.Wallet{}).Error
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	return tx.Where("id IN ?", ids).Delete(&model.Wallet{}).Error
}

func (Wallet) Add(ctx *gin.Context) {
	var req wAddReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	kind := normalizeWalletKind(req.Kind)
	baseType := normalizeWalletBase(req.Base)
	tradeTypes, err := resolveWalletTradeTypes(kind, baseType, req.TradeType, req.TradeTypes)
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	rows, err := buildWalletRows(model.NewWalletGroupID(), req.Name, req.Remark, req.Address, model.WaStatusEnable, req.OtherNotify, kind, baseType, tradeTypes)
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	if err := model.Db.Create(&rows).Error; err != nil {
		base.Error(ctx, err)
		return
	}

	base.Ok(ctx, "success")
}

func (Wallet) List(ctx *gin.Context) {
	var req wListReq
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	var rows []model.Wallet
	db := model.Db.Order("id " + req.Sort)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Address != "" {
		db = db.Where("address LIKE ?", "%"+req.Address+"%")
	}
	if req.Trade != "" {
		db = db.Where("trade_type LIKE ?", "%"+req.Trade+"%")
	}

	if err := db.Find(&rows).Error; err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	groups := model.GroupWallets(rows)
	total := int64(len(groups))
	start := (req.Page - 1) * req.Size
	if start < 0 {
		start = 0
	}
	if start > len(groups) {
		start = len(groups)
	}
	end := start + req.Size
	if end > len(groups) {
		end = len(groups)
	}

	base.Response(ctx, 200, groups[start:end], total)
}

func (Wallet) Mod(ctx *gin.Context) {
	var req wModReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	kind := normalizeWalletKind(req.Kind)
	baseType := normalizeWalletBase(req.Base)
	tradeTypes, err := resolveWalletTradeTypes(kind, baseType, req.TradeType, req.TradeTypes)
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	rows, err := loadWalletGroupRows(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			base.BadRequest(ctx, "钱包不存在")
			return
		}
		base.Error(ctx, err)
		return
	}

	storedGroupID := strings.TrimSpace(rows[0].GroupID)
	targetGroupID := storedGroupID
	if targetGroupID == "" {
		targetGroupID = model.NewWalletGroupID()
	}

	newRows, err := buildWalletRows(targetGroupID, req.Name, req.Remark, req.Address, req.Status, req.OtherNotify, kind, baseType, tradeTypes)
	if err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	if err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := deleteWalletRows(tx, rows); err != nil {
			return err
		}
		return tx.Create(&newRows).Error
	}); err != nil {
		base.Error(ctx, err)
		return
	}

	base.Ok(ctx, "修改成功")
}

func (Wallet) Del(ctx *gin.Context) {
	var req base.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())
		return
	}

	rows, err := loadWalletGroupRows(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			base.BadRequest(ctx, "钱包不存在")
			return
		}
		base.Error(ctx, err)
		return
	}

	if err := deleteWalletRows(model.Db, rows); err != nil {
		base.Error(ctx, err)
		return
	}

	base.Ok(ctx, "删除成功")
}
