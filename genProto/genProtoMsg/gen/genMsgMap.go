package gen

import (
	"bufio"
	"fmt"
	"genProto/genProtoMsg/conf"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

type IdMapUnit struct {
	MsgId   uint16
	MsgName string
}

type IdMapData struct {
	Package string
	List    []IdMapUnit
}

func parseLine(n string) string {
	str := strings.ReplaceAll(n, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, ";", "")
	str = strings.ReplaceAll(str, `"`, "")
	return str
}

func splits(s string) []string {
	return strings.Split(s, "=")
}

func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func filter[T comparable](data []T, f func(T) bool) []T {
	var result []T
	for _, v := range data {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func isNotEmpty(s string) bool {
	return s != ""
}

func loadIdMapData() IdMapData {
	list := make([]IdMapUnit, 0)

	data := IdMapData{}
	c := conf.GetConf()
	nameSlice := strings.Split(c.Source, "/")
	secondLast := nameSlice[len(nameSlice)-2]
	// 取出第一个字符并转换为rune
	firstCharRune := []rune(secondLast)[0]
	// 转换为小写
	firstCharLower := unicode.ToLower(firstCharRune)
	protoName := fmt.Sprintf("%v_msgId.proto", string(firstCharLower))

	file, err := os.Open(c.Source + protoName)
	if err != nil {
		fmt.Println("open file err:", err)
		return data
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		n := scanner.Text()

		if strings.Contains(n, `go_package`) {
			n = strings.ReplaceAll(n, `"`, "")
			packageSlice := strings.Split(n, ";")
			packageSlice = filter(packageSlice, isNotEmpty)
			data.Package = packageSlice[len(packageSlice)-1]
			continue
		}

		n = parseLine(n)
		if len(n) == 0 {
			continue
		}

		subs := []string{"syntax", "enum", "}", "//", "package"}
		if ContainsAll(n, subs...) {
			continue
		}

		str := splits(n)
		if len(str) == 0 {
			continue
		}

		if len(str) != 2 {
			log.Fatal("msgId err. data:", n)
		}

		msgId, _ := strconv.Atoi(str[1])
		unit := IdMapUnit{
			MsgId:   uint16(msgId),
			MsgName: str[0],
		}
		list = append(list, unit)
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	data.List = list
	return data
}

func replaceName(name string) string {
	str := strings.ReplaceAll(name, "ID_", "")
	str = strings.ReplaceAll(str, `_`, "")
	return str
}

func rangeInit(list []IdMapUnit) string {
	str := ""
	for _, v := range list {
		str += fmt.Sprintf("\tmsgProtoTypes[Msg%sId]=reflect.TypeOf((*%s)(nil)).Elem()\n", v.MsgName, replaceName(v.MsgName))
	}

	for _, v := range list {
		str += fmt.Sprintf("\tmsgNames[Msg%sId]=\"%s\"\n", v.MsgName, replaceName(v.MsgName))
	}
	str = str[:len(str)-1]

	return str
}

func rangeConst(list []IdMapUnit) string {
	str := ""

	for _, v := range list {
		str += fmt.Sprintf("\tMsg%sId=%d\n", v.MsgName, v.MsgId)
	}
	str = str[:len(str)-1]

	return str
}

func rangeFormType(list []IdMapUnit) string {
	str := ""

	for _, v := range list {
		str += fmt.Sprintf("\tcase *%s:\n", replaceName(v.MsgName))
		str += fmt.Sprintf("\t\treturn %d\n", v.MsgId)
	}
	str = str[:len(str)-1]

	return str
}

func GenMsgMap(fileName string) {
	funcMap := template.FuncMap{
		"fnToTitle":     strings.ToTitle,
		"rangeInit":     rangeInit,
		"rangeConst":    rangeConst,
		"rangeFormType": rangeFormType,
	}

	const templateText = `package {{.Package}}
import (
	"reflect"
)

var msgProtoTypes = make(map[uint16]reflect.Type)
var msgNames = make(map[uint16]string)

func init() {
{{.List | rangeInit}}
}

func GetMsgProtoType(key uint16) reflect.Type {
	return msgProtoTypes[key]
}

func GetMsgName(key uint16) string {
	return msgNames[key]
}

const (
{{.List | rangeConst}}
)

func GetMsgIdFromType(i interface{}) uint16 {
	switch i.(type) {
{{.List | rangeFormType}}
	default:
		return 0
	}
}
`
	tpl := template.New("example").Funcs(funcMap)
	tpl, _ = tpl.Parse(templateText)

	d := loadIdMapData()
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("openFile err", err)
		return
	}

	tpl.Execute(file, d)
}
