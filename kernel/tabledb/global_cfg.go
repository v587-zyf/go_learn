package tabledb

type GlobalBaseCfg struct {
	Id    int    `col:"id" client:"id"`
	Type  string `col:"clinetType" client:"clinetType"`
	Name  string `col:"name" client:"name"`
	Value string `col:"value" client:"value"`
}
