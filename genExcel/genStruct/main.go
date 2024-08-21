package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	codePackage      = flag.String("package", "tdb", "code package")
	savePath         = flag.String("savePath", "", "Path to save the makefile")
	readPath         = flag.String("readPath", "", "The path of reading Excel")
	allType          = flag.String("allType", "", "Specified field type")
	generateLanguage = flag.Bool("l", false, "Path to save the makefile")
	generateClient   = flag.Bool("c", false, "Path to save the makefile")
	withoutExcel     = flag.String("withoutExcel", "", "过滤表")
)

func main() {
	flag.Parse()
	if *savePath == "" || *readPath == "" {
		fmt.Printf("savePath, readPath or allType is nil\n")
		return
	}

	//fmt.Println("----------------")
	//fmt.Println("开始生成语言包")
	//
	//GenText(*readPath, *savePath)
	//
	//fmt.Println("结束生成语言包")
	//fmt.Println("----------------")
	//if *generateLanguage {
	//	return
	//}

	if *allType == "" {
		fmt.Println("allType 不能为空")
		return
	}

	fmt.Println("----------------")
	fmt.Println("开始表结构生产")

	withoutExcelSlice := strings.Split(*withoutExcel, ",")
	withoutExcelMap := make(map[string]bool)
	for _, v := range withoutExcelSlice {
		s := strings.TrimSpace(v)
		if len(v) > 0 {
			withoutExcelMap[s+".xlsx"] = true
		}
	}

	gt := Generate{}
	err := gt.ReadExcel(*readPath, *savePath, *allType, *generateClient, withoutExcelMap)
	if err != nil {
		fmt.Printf("something err:%v\n", err)
		return
	}

	fmt.Println("结束表结构生产")
	fmt.Println("----------------")
}
