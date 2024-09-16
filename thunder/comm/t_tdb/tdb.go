package t_tdb

import (
	"encoding/gob"
	"github.com/v587-zyf/gc/iface"
	c "github.com/v587-zyf/gc/tableDb"
)

type Tdb struct {
	*c.TableDb
	*TableBase
	TableDbPath string

	CardMap   map[int][]*CardCardCfg
	CardLvMap map[int]map[int]*CardLvCardLvCfg
	ShopLvMap map[int]map[int]*ShopShopCfg
	LvSlice   []int
}

var t_tdb *Tdb

func init() {
	t_tdb = &Tdb{
		TableBase: &TableBase{},
		TableDb: &c.TableDb{
			FileModTime: make(map[string]int64),
			InitConf:    &InitConf{},
		},
	}
	//tdb.TableDb.FileModTime = make(map[string]int64)
	//tdb.TableDb.InitConf = &InitConf{}
}

func Init(tableDbPath string) error {
	t_tdb.Init(tableDbPath)

	return t_tdb.Load(t_tdb)
}

func (t *Tdb) Load(tdb iface.ITableDb) (err error) {
	t.TableDb.FileInfos = fileInfos

	otherReg := []any{
		InitConf{},
	}
	for _, a := range otherReg {
		gob.Register(a)
	}
	//gob.Register(TextErrorTextCfg{})
	for _, info := range t.FileInfos {
		for _, i := range info.SheetInfos {
			gob.Register(i.ObjPropType)
		}
	}

	err = t.TableDb.Load(tdb)
	if err != nil {
		return err
	}

	err = t.CheckConf(t.InitConf)

	t.Patch()

	return
}
func (t *Tdb) CheckConf(c any) (err error) {
	return t.TableDb.CheckConf(t_tdb)
}

func Conf() *InitConf {
	if conf, ok := t_tdb.InitConf.(InitConf); ok {
		return &conf
	} else {
		return t_tdb.InitConf.(*InitConf)
	}
}

func Get() *Tdb {
	return t_tdb
}
