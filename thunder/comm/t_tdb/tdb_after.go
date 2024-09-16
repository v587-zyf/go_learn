package t_tdb

func (t *Tdb) Patch() {
	t.genCard()
	t.genCardLv()
	t.genShopLv()
	t.genLvSlice()
}

func (t *Tdb) genCard() {
	t.CardMap = make(map[int][]*CardCardCfg)
	for _, cfg := range t.CardCardCfgs {
		if t.CardMap[cfg.Map_id] == nil {
			t.CardMap[cfg.Map_id] = make([]*CardCardCfg, 0)
		}
		t.CardMap[cfg.Map_id] = append(t.CardMap[cfg.Map_id], cfg)
	}
}
func (t *Tdb) genCardLv() {
	t.CardLvMap = make(map[int]map[int]*CardLvCardLvCfg)
	for _, cfg := range t.CardLvCardLvCfgs {
		if t.CardLvMap[cfg.Card_id] == nil {
			t.CardLvMap[cfg.Card_id] = make(map[int]*CardLvCardLvCfg)
		}
		t.CardLvMap[cfg.Card_id][cfg.Card_lv] = cfg
	}
}
func (t *Tdb) genShopLv() {
	t.ShopLvMap = make(map[int]map[int]*ShopShopCfg)
	for _, cfg := range t.ShopShopCfgs {
		if t.ShopLvMap[cfg.Type] == nil {
			t.ShopLvMap[cfg.Type] = make(map[int]*ShopShopCfg)
		}
		t.ShopLvMap[cfg.Type][cfg.Buy_num] = cfg
	}
}
func (t *Tdb) genLvSlice() {
	i := 0
	t.LvSlice = make([]int, len(t.UserUserCfgs))
	for _, cfg := range t.UserUserCfgs {
		t.LvSlice[i] = cfg.Id
		i++
	}
}
