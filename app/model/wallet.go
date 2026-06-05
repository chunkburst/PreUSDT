package model

import (
	"sort"
	"strings"
	"time"

	"github.com/chunkburst/PreUSDT/app/utils"
)

const (
	WaStatusEnable  uint8 = 1
	WaStatusDisable uint8 = 0
	WaOtherEnable   uint8 = 1
	WaOtherDisable  uint8 = 0
)

type Wallet struct {
	Id
	Name        string     `gorm:"column:name;type:varchar(32);not null;default:-';comment:名称" json:"name"`
	Status      uint8      `gorm:"column:status;not null;default:1;comment:地址状态" json:"status"`
	Address     string     `gorm:"column:address;type:varchar(128);not null;index;comment:钱包地址" json:"address"`
	MatchAddr   string     `gorm:"column:match_addr;type:varchar(128);not null;uniqueIndex:idx_address;comment:匹配地址" json:"match_addr"`
	TradeType   string     `gorm:"column:trade_type;type:varchar(20);not null;uniqueIndex:idx_address;comment:交易类型" json:"trade_type"`
	OtherNotify uint8      `gorm:"column:other_notify;not null;default:0;comment:其它通知" json:"other_notify"`
	Remark      string     `gorm:"column:remark;type:varchar(255);not null;default:'';comment:备注" json:"remark"`
	GroupID     string     `gorm:"column:group_id;type:varchar(32);not null;default:'';index;comment:钱包分组ID" json:"group_id"`
	Kind        WalletKind `gorm:"column:kind;type:varchar(16);not null;default:'single';comment:钱包类型" json:"kind"`
	Base        WalletBase `gorm:"column:base;type:varchar(16);not null;default:'generic';comment:基础分类" json:"base"`
	AutoTimeAt
}

type WalletGroup struct {
	ID          int64      `json:"id"`
	GroupID     string     `json:"group_id"`
	Name        string     `json:"name"`
	Status      uint8      `json:"status"`
	Address     string     `json:"address"`
	Remark      string     `json:"remark"`
	OtherNotify uint8      `json:"other_notify"`
	Kind        WalletKind `json:"kind"`
	Base        WalletBase `json:"base"`
	TradeType   string     `json:"trade_type"`
	TradeTypes  []string   `json:"trade_types"`
	CreatedAt   *Datetime  `json:"created_at,omitempty"`
	UpdatedAt   *Datetime  `json:"updated_at,omitempty"`
}

type WalletBaseOption struct {
	Value      WalletBase `json:"value"`
	Label      string     `json:"label"`
	TradeTypes []string   `json:"trade_types"`
}

var walletBaseTradeTypes = map[WalletBase][]TradeType{
	WalletBaseEvm: {
		EthereumEth,
		BscBnb,
		UsdtErc20,
		UsdcErc20,
		UsdtBep20,
		UsdcBep20,
		UsdtPolygon,
		UsdcPolygon,
		UsdtArbitrum,
		UsdcArbitrum,
		UsdtXlayer,
		UsdcXlayer,
		UsdcBase,
	},
}

var walletBaseLabels = map[WalletBase]string{
	WalletBaseGeneric: "通用地址",
	WalletBaseEvm:     "EVM 兼容地址",
}

func (wa *Wallet) TableName() string {
	return "bep_wallet"
}

func (wa *Wallet) SetStatus(status uint8) {
	wa.Status = status
	Db.Save(wa)
}

func (wa *Wallet) IsValid() bool {
	tradeType := TradeType(wa.TradeType)

	if tradeType == TronTrx || tradeType == UsdtTrc20 || tradeType == UsdcTrc20 {
		return utils.IsValidTronAddress(wa.Address)
	}
	if tradeType == UsdtSolana || tradeType == UsdcSolana {
		return utils.IsValidSolanaAddress(wa.Address)
	}
	if tradeType == UsdtAptos || tradeType == UsdcAptos {
		return utils.IsValidAptosAddress(wa.Address)
	}

	return utils.IsValidEvmAddress(wa.Address)
}

func (wa *Wallet) SetOtherNotify(notify uint8) {
	wa.OtherNotify = notify
	Db.Save(wa)
}

func (wa *Wallet) Delete() {
	Db.Delete(wa)
}

func (wa *Wallet) GetTokenContract() string {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {
		return c.Contract
	}

	return ""
}

func (wa *Wallet) GetTokenDecimals() int32 {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {
		return c.Decimal
	}

	return -18
}

func (wa *Wallet) GetNetwork() Network {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {
		return c.Network
	}

	return ""
}

func GetAvailableAddress(t TradeType) []string {
	var rows []Wallet
	Db.Where("trade_type = ? and status = ?", t, WaStatusEnable).Find(&rows)

	wallets := make([]string, 0, len(rows))
	for _, w := range rows {
		wallets = append(wallets, w.MatchAddr)
	}

	return wallets
}

func (wa *Wallet) NormalizeAddress() {
	wa.Address = strings.TrimSpace(wa.Address)
	wa.MatchAddr = wa.Address
	if !AddrCaseSens(TradeType(wa.TradeType)) {
		wa.MatchAddr = strings.ToLower(wa.MatchAddr)
	}
}

func (wa *Wallet) FillGrouping() {
	if wa.GroupID == "" {
		wa.GroupID = wa.TradeType + ":" + wa.MatchAddr
	}
	if wa.Kind == "" {
		wa.Kind = WalletKindSingle
	}
	if wa.Base == "" {
		wa.Base = WalletBaseGeneric
	}
}

func NewWalletGroupID() string {
	id, err := utils.GenerateTradeId()
	if err != nil {
		return utils.Md5String(utils.StrSha256(time.Now().String()))[:18]
	}

	return id
}

func GetWalletBaseTradeTypes(base WalletBase) []TradeType {
	list, ok := walletBaseTradeTypes[base]
	if !ok {
		return []TradeType{}
	}

	copied := make([]TradeType, len(list))
	copy(copied, list)
	return copied
}

func IsTradeTypeAllowedForWalletBase(base WalletBase, tradeType TradeType) bool {
	for _, item := range GetWalletBaseTradeTypes(base) {
		if item == tradeType {
			return true
		}
	}

	return false
}

func GetWalletBaseOptions() []WalletBaseOption {
	options := make([]WalletBaseOption, 0, len(walletBaseTradeTypes))
	for base, tradeTypes := range walletBaseTradeTypes {
		items := make([]string, 0, len(tradeTypes))
		for _, tradeType := range tradeTypes {
			items = append(items, string(tradeType))
		}
		sort.Strings(items)
		options = append(options, WalletBaseOption{
			Value:      base,
			Label:      walletBaseLabels[base],
			TradeTypes: items,
		})
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Value < options[j].Value
	})

	return options
}

func GroupWallets(rows []Wallet) []WalletGroup {
	groups := make(map[string]*WalletGroup)
	order := make([]string, 0)

	for _, row := range rows {
		row.FillGrouping()
		item, ok := groups[row.GroupID]
		if !ok {
			tradeType := row.TradeType
			if row.Kind == WalletKindMulti {
				tradeType = ""
			}
			item = &WalletGroup{
				ID:          row.ID,
				GroupID:     row.GroupID,
				Name:        row.Name,
				Status:      row.Status,
				Address:     row.Address,
				Remark:      row.Remark,
				OtherNotify: row.OtherNotify,
				Kind:        row.Kind,
				Base:        row.Base,
				TradeType:   tradeType,
				TradeTypes:  make([]string, 0, 1),
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
			groups[row.GroupID] = item
			order = append(order, row.GroupID)
		}

		item.TradeTypes = append(item.TradeTypes, row.TradeType)
		if row.ID < item.ID {
			item.ID = row.ID
		}
		if row.Status < item.Status {
			item.Status = row.Status
		}
		if row.OtherNotify > item.OtherNotify {
			item.OtherNotify = row.OtherNotify
		}
		if row.UpdatedAt != nil && (item.UpdatedAt == nil || item.UpdatedAt.Before(row.UpdatedAt.Time())) {
			item.UpdatedAt = row.UpdatedAt
		}
		if row.CreatedAt != nil && (item.CreatedAt == nil || row.CreatedAt.Before(item.CreatedAt.Time())) {
			item.CreatedAt = row.CreatedAt
		}
	}

	list := make([]WalletGroup, 0, len(order))
	for _, key := range order {
		item := groups[key]
		sort.Strings(item.TradeTypes)
		list = append(list, *item)
	}

	return list
}
