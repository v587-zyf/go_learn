package db

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type Model interface {
	SetDbMap(dbMap *gorp.DbMap) // 设置数据库集合
	DbMap() *gorp.DbMap         // 获取数据库集合
	SetDb(db *sql.DB)           // 设置数据库
	Db() *sql.DB                // 获取数据库
}

type modelMapItem struct {
	model  Model
	initer func(dbMap *gorp.DbMap)
}

var (
	cache           *memcache.Client //缓存信息
	cacheExpiration int32            //缓存到期时间
	enableDbTrace   bool             //缓存跟踪开关

	modelMap  = make(map[string][]modelMapItem) //存所有数据库信息
	callbacks = make([]func(), 0)
	inited    = false
)

// 公共数据库结构体
type CommonModel struct {
	dbMap *gorp.DbMap
	db    *sql.DB
}

func init() {
	// 是否开启跟踪
	flag.BoolVar(&enableDbTrace, "db_trace", false, "whether enable gorp db trace")
}

func (this *CommonModel) SetMap(dbMap *gorp.DbMap) {
	this.dbMap = dbMap
}
func (this *CommonModel) DbMap() *gorp.DbMap {
	return this.dbMap
}
func (this *CommonModel) SetDb(db *sql.DB) {
	this.db = db
}
func (this *CommonModel) Db() *sql.DB {
	return this.db
}

// 注册新数据库文件
func Register(dbKey string, model Model, initer func(dbMap *gorp.DbMap)) {
	// 判断是否已有
	mapItems, ok := modelMap[dbKey]
	// 没有加入到内存
	if ok {
		for _, mi := range mapItems {
			if model == mi.model {
				return
			}
			mapItems = append(modelMap[dbKey], modelMapItem{model: model, initer: initer})
		}
	} else {
		mapItems := make([]modelMapItem, 0, 5)
		mapItems = append(mapItems, modelMapItem{model: model, initer: initer})
		modelMap[dbKey] = mapItems
	}
}

type DbConnectionStringGetter interface {
	GetDbConnectionString(dbKey string) (string, int, int)
}

type DbLogger struct {
}

func (this *DbLogger) Printf(format string, v ...interface{}) {
	fmt.Errorf(format, v...)
}

// 初始化数据库
func InitDb(connStrGetter DbConnectionStringGetter, ormDefaultDbKey string, tableAutoCheckDbKey, columnCheck []string) error {
	for dbKey := range modelMap {
		dbUrl, maxIdle, maxOpenCon := connStrGetter.GetDbConnectionString(dbKey)
		if dbUrl == "" {
			continue
		}
		// 检查是否需要表自动检查
		tableAutoCheck := false
		if tableAutoCheckDbKey != nil {
			for _, v := range tableAutoCheckDbKey {
				if dbKey == v {
					tableAutoCheck = true
					break
				}
			}
		}
		columnAutoCheck := false
		// 检查是否自动检查表字段
		if columnCheck != nil {
			for _, v := range columnCheck {
				if dbKey == v {
					columnAutoCheck = true
					break
				}
			}
		}
		err := dbInit(dbUrl, dbKey, dbKey == ormDefaultDbKey, tableAutoCheck, columnAutoCheck, maxIdle, maxOpenCon)
		if err != nil {
			return err
		}
	}

	inited = true
	for _, callback := range callbacks {
		callback()
	}

	return nil
}

// 初始化数据库详细操作
func dbInit(dbUrl, dbKey string, ormDefaultDb, tableAutoCheck, columnAutoCheck bool, maxIdle, maxOpenCon int) error {
	if dbUrl == "" {
		return fmt.Errorf("init db:%v error,db has not config", dbKey)
	}
	// 连接
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return err
	}
	// 测试连接
	err = db.Ping()
	if err != nil {
		return err
	}
	// 最大空闲连接数
	if maxIdle <= 0 {
		maxIdle = 3
	}
	// 最大打开连接数
	if maxOpenCon <= 0 {
		maxOpenCon = 5
	}
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpenCon)

	// 设置连接最大生命周期
	//db.SetConnMaxLifetime(110)

	// 是否跟踪错误
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	if enableDbTrace {
		dbMap.TraceOn("[gorp]", &DbLogger{})
	}

	// 设置所有库
	mapItems, ok := modelMap[dbKey]
	if ok {
		for _, mi := range mapItems {
			mi.model.SetDbMap(dbMap)
			mi.model.SetDb(db)
			mi.initer(dbMap)
		}
	}

	// 是否注册为 dbOrm 默认数据库
	if ormDefaultDb {
		err := orm.AddAliasWthDB("default", "mysql", db)
		if err != nil {
			return err
		}
	}

	// 表自动检测结构
	if tableAutoCheck {
		err := orm.AddAliasWthDB(dbKey, "mysql", db)
		if err != nil {
			return err
		}
		// 创建table
		err = orm.RunSyncdb(dbKey, false, true)
		if err != nil {
			return err
		}
	}

	// 自动检测列
	//if columnAutoCheck {
	//	dbMap
	//}

	return nil
}
