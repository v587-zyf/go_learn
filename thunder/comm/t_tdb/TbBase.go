package t_tdb

import c "github.com/v587-zyf/gc/tableDb"

var fileInfos = []c.FileInfo{
	{"globals.xlsx", []c.SheetInfo{
		{"global", c.LoadGlobalConf, c.GlobalBaseCfg{}},
	}},

	{"card.xlsx", []c.SheetInfo{
		{SheetName: "card", Initer: c.MapLoader("CardCardCfgs", "Id"), ObjPropType: CardCardCfg{}},
	}},
	{"card_lv.xlsx", []c.SheetInfo{
		{SheetName: "card_lv", Initer: c.MapLoader("CardLvCardLvCfgs", "Id"), ObjPropType: CardLvCardLvCfg{}},
	}},
	{"invite_lv.xlsx", []c.SheetInfo{
		{SheetName: "invite_lv", Initer: c.MapLoader("InviteLvInviteLvCfgs", "Id"), ObjPropType: InviteLvInviteLvCfg{}},
	}},
	{"invite_num.xlsx", []c.SheetInfo{
		{SheetName: "invite_num", Initer: c.MapLoader("InviteNumInviteNumCfgs", "Id"), ObjPropType: InviteNumInviteNumCfg{}},
	}},
	{"map.xlsx", []c.SheetInfo{
		{SheetName: "map", Initer: c.MapLoader("MapMapCfgs", "Id"), ObjPropType: MapMapCfg{}},
	}},
	{"pay.xlsx", []c.SheetInfo{
		{SheetName: "pay", Initer: c.MapLoader("PayPayCfgs", "Id"), ObjPropType: PayPayCfg{}},
	}},
	{"shop.xlsx", []c.SheetInfo{
		{SheetName: "shop", Initer: c.MapLoader("ShopShopCfgs", "Id"), ObjPropType: ShopShopCfg{}},
	}},
	{"treasure.xlsx", []c.SheetInfo{
		{SheetName: "treasure", Initer: c.MapLoader("TreasureTreasureCfgs", "Id"), ObjPropType: TreasureTreasureCfg{}},
	}},
	{"user.xlsx", []c.SheetInfo{
		{SheetName: "user", Initer: c.MapLoader("UserUserCfgs", "Id"), ObjPropType: UserUserCfg{}},
	}},
}

type TableBase struct {
	// NOTE 关于client的配置：
	// client:对象名,对象类型　，对象名要小写．
	// mapKey 即对应的我们的结构里的key, 要看具体的型中key是什么　，一段是大写的

	CardCardCfgs           map[int]*CardCardCfg
	CardLvCardLvCfgs       map[int]*CardLvCardLvCfg
	InviteLvInviteLvCfgs   map[int]*InviteLvInviteLvCfg
	InviteNumInviteNumCfgs map[int]*InviteNumInviteNumCfg
	MapMapCfgs             map[int]*MapMapCfg
	PayPayCfgs             map[int]*PayPayCfg
	ShopShopCfgs           map[int]*ShopShopCfg
	TreasureTreasureCfgs   map[int]*TreasureTreasureCfg
	UserUserCfgs           map[int]*UserUserCfg
}

func GetCardCardCfg(Id int) *CardCardCfg {
	return t_tdb.CardCardCfgs[Id]
}

func RangCardCardCfgs(f func(conf *CardCardCfg) bool) {
	for _, v := range t_tdb.CardCardCfgs {
		if !f(v) {
			return
		}
	}
}

func GetCardLvCardLvCfg(Id int) *CardLvCardLvCfg {
	return t_tdb.CardLvCardLvCfgs[Id]
}

func RangCardLvCardLvCfgs(f func(conf *CardLvCardLvCfg) bool) {
	for _, v := range t_tdb.CardLvCardLvCfgs {
		if !f(v) {
			return
		}
	}
}

func GetInviteLvInviteLvCfg(Id int) *InviteLvInviteLvCfg {
	return t_tdb.InviteLvInviteLvCfgs[Id]
}

func RangInviteLvInviteLvCfgs(f func(conf *InviteLvInviteLvCfg) bool) {
	for _, v := range t_tdb.InviteLvInviteLvCfgs {
		if !f(v) {
			return
		}
	}
}

func GetInviteNumInviteNumCfg(Id int) *InviteNumInviteNumCfg {
	return t_tdb.InviteNumInviteNumCfgs[Id]
}

func RangInviteNumInviteNumCfgs(f func(conf *InviteNumInviteNumCfg) bool) {
	for _, v := range t_tdb.InviteNumInviteNumCfgs {
		if !f(v) {
			return
		}
	}
}

func GetMapMapCfg(Id int) *MapMapCfg {
	return t_tdb.MapMapCfgs[Id]
}

func RangMapMapCfgs(f func(conf *MapMapCfg) bool) {
	for _, v := range t_tdb.MapMapCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPayPayCfg(Id int) *PayPayCfg {
	return t_tdb.PayPayCfgs[Id]
}

func RangPayPayCfgs(f func(conf *PayPayCfg) bool) {
	for _, v := range t_tdb.PayPayCfgs {
		if !f(v) {
			return
		}
	}
}

func GetShopShopCfg(Id int) *ShopShopCfg {
	return t_tdb.ShopShopCfgs[Id]
}

func RangShopShopCfgs(f func(conf *ShopShopCfg) bool) {
	for _, v := range t_tdb.ShopShopCfgs {
		if !f(v) {
			return
		}
	}
}

func GetTreasureTreasureCfg(Id int) *TreasureTreasureCfg {
	return t_tdb.TreasureTreasureCfgs[Id]
}

func RangTreasureTreasureCfgs(f func(conf *TreasureTreasureCfg) bool) {
	for _, v := range t_tdb.TreasureTreasureCfgs {
		if !f(v) {
			return
		}
	}
}

func GetUserUserCfg(Id int) *UserUserCfg {
	return t_tdb.UserUserCfgs[Id]
}

func RangUserUserCfgs(f func(conf *UserUserCfg) bool) {
	for _, v := range t_tdb.UserUserCfgs {
		if !f(v) {
			return
		}
	}
}

type CardCardCfg struct {
	Id     int `col:"id" client:"id"`         // 卡牌id
	Map_id int `col:"map_id" client:"map_id"` // 获取最低地图等级
	Weight int `col:"weight" client:"weight"` // 权重
}

type CardLvCardLvCfg struct {
	Id              int `col:"id" client:"id"`                           // id
	Card_id         int `col:"card_id" client:"card_id"`                 // 卡牌id
	Card_lv         int `col:"card_lv" client:"card_lv"`                 // 卡牌等级
	Consume_gold    int `col:"consume_gold" client:"consume_gold"`       // 升级金币
	Consume_diamond int `col:"consume_diamond" client:"consume_diamond"` // 升级钻石
	Income          int `col:"income" client:"income"`                   // 每小时收益
}

type InviteLvInviteLvCfg struct {
	Id     int `col:"id" client:"id"`         // 玩家等级
	Normal int `col:"normal" client:"normal"` // 普通玩家奖励
	Tg_vip int `col:"tg_vip" client:"tg_vip"` // TG会员奖励
}

type InviteNumInviteNumCfg struct {
	Id      int `col:"id" client:"id"`           // id
	Num     int `col:"num" client:"num"`         // 分享人数
	Card_id int `col:"card_id" client:"card_id"` // 获取卡牌id
}

type MapMapCfg struct {
	Id             int `col:"id" client:"id"`                         // 地图等级
	Long           int `col:"long" client:"long"`                     // 长
	Width          int `col:"width" client:"width"`                   // 宽
	Reset_strength int `col:"reset_strength" client:"reset_strength"` // 重置地图消耗体力
}

type PayPayCfg struct {
	Id      int    `col:"id" client:"id"`           // id
	Money   string `col:"money" client:"money"`     // 充值金额
	Diamond int    `col:"diamond" client:"diamond"` // 获得钻石
}

type ShopShopCfg struct {
	Id      int `col:"id" client:"id"`           // id
	Type    int `col:"type" client:"type"`       // 商品类型
	Buy_num int `col:"buy_num" client:"buy_num"` // 购买次数
	Diamond int `col:"diamond" client:"diamond"` // 消耗钻石
}

type TreasureTreasureCfg struct {
	Id   int `col:"id" client:"id"`     // id
	Fold int `col:"fold" client:"fold"` // 宝箱金币倍率
	Rate int `col:"rate" client:"rate"` // 概率(万分比)
}

type UserUserCfg struct {
	Id           int `col:"id" client:"id"`                     // 玩家等级
	Consume      int `col:"consume" client:"consume"`           // 升级金币
	Map_id       int `col:"map_id" client:"map_id"`             // 地图等级
	Dig_strength int `col:"dig_strength" client:"dig_strength"` // 挖墙消耗体力
	Strength_max int `col:"strength_max" client:"strength_max"` // 体力上限
	Gacha_num    int `col:"gacha_num" client:"gacha_num"`       // 抽卡次数
}

type InitConf struct {
	Thunder_rate          int `conf:"thunder_rate" default:"1560"`           //地雷战比(万分比)
	Treasure_rate         int `conf:"treasure_rate" default:"500"`           //宝箱占比(万分比)
	Revive_strength       int `conf:"revive_strength" default:"20"`          //复活消耗体力
	Add_strengthen_second int `conf:"add_strengthen_second" default:"10"`    //恢复1次体力秒数
	Add_strengthen_num    int `conf:"add_strengthen_num" default:"3"`        //1次恢复体力数量
	Init_card_id          int `conf:"init_card_id" default:"1"`              //初始卡牌id
	Normal_invite_diamond int `conf:"normal_invite_diamond" default:"1"`     //普通玩家邀请钻石奖励数量
	Vip_invite_diamond    int `conf:"vip_invite_diamond" default:"10"`       //TG会员邀请钻石奖励数量
	Buy_extra_strength    int `conf:"buy_extra_strength" default:"500"`      //每次购买额外体力
	Free_hasten_cd_second int `conf:"free_hasten_cd_second" default:"7200"`  //免费加速间隔秒数
	Free_hasten_second    int `conf:"free_hasten_second" default:"900"`      //免费加速加速秒数
	Link_hasten_cd_second int `conf:"link_hasten_cd_second" default:"43200"` //链上签到间隔秒数
	Link_hasten_second    int `conf:"link_hasten_second" default:"5400"`     //链上签到加速秒数
	Pay_hasten_diamond    int `conf:"pay_hasten_diamond" default:"16"`       //付费加速消耗钻石
	Pay_hasten_second     int `conf:"pay_hasten_second" default:"14400"`     //付费加速秒数
	Hasten_max_second     int `conf:"hasten_max_second" default:"43200"`     //加速时长大于几秒不可再加速

}
