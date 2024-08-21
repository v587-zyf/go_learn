package db

type BaseTableInterface interface {
	TableName() string       // 表名
	TableEngine() string     // 引擎
	TableEncode() string     // 编码
	TableComment() string    // 注释
	TableIndex() [][]string  // 单字段索引
	TableUnique() [][]string // 多字段索引
}

type BaseTable struct {
}

func (this *BaseTable) TableName() string {
	return "defName"
}

func (this *BaseTable) TableEngine() string {
	return "Innodb"
}

func (this *BaseTable) TableEncode() string {
	return "utf8"
}

func (this *BaseTable) TableComment() string {
	return ""
}

func (this *BaseTable) TableIndex() [][]string {
	return [][]string{}
}

func (this *BaseTable) TableUnique() [][]string {
	return [][]string{}
}
