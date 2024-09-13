package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
)

func (g *Generate) genGlobalConf(readPath string) string {
	wb, err := xlsx.OpenFile(readPath + "\\globals.xlsx")
	if err != nil {
		fmt.Println(fmt.Sprintf("ReadExcel|xlsx.OpenFile is err :%v", err))
	}
	// 遍历工作表
	var confFiled string
	for _, sheet := range wb.Sheets {
		if hasChineseOrDefault(sheet.Name) {
			continue
		}

		typeCell := -1
		nameCell := -1
		valueCell := -1
		descCell := -1
		maxCell := -1
		for k, v := range sheet.Rows[0].Cells {
			//fmt.Printf("k:%v val:%v\n", k, v.Value)
			if v.Value == "clientType" {
				typeCell = k
			} else if v.Value == "name" {
				nameCell = k
			} else if v.Value == "value" {
				valueCell = k
			} else if v.Value == "desc" {
				descCell = k
			} else {
				continue
			}
		}
		if maxCell < typeCell {
			maxCell = typeCell
		} else if maxCell < nameCell {
			maxCell = nameCell
		} else if maxCell < valueCell {
			maxCell = valueCell
		} else if maxCell < descCell {
			maxCell = descCell
		}
		if typeCell == -1 || nameCell == -1 || valueCell == -1 || descCell == -1 {
			panic(fmt.Sprintf("global typeCell:%v nameCell:%v valueCell:%v descCell:%v", typeCell, nameCell, valueCell, descCell))
		}

		for k, v := range sheet.Rows {
			if k < lineNumber {
				continue
			}
			if len(v.Cells) == 0 || len(v.Cells) < maxCell {
				break
			}
			//fmt.Printf("name:%v type:%v value:%v desc:%v\n", nameCell, typeCell, valueCell, descCell)
			//fmt.Printf("name:%v type:%v value:%v desc:%v\n", v.Cells[nameCell].Value, v.Cells[typeCell].Value, v.Cells[valueCell].Value, v.Cells[descCell].Value)
			//fmt.Println(nameCell, "---", len(v.Cells))
			confFiled += fmt.Sprintf(FIELD_TEMP,
				FirstRuneToUpper(v.Cells[nameCell].Value),
				g.CheckType(v.Cells[typeCell].Value, "global"),
				strings.TrimSpace(v.Cells[nameCell].Value),
				strings.TrimSpace(v.Cells[valueCell].Value),
				strings.TrimSpace(v.Cells[descCell].Value))
		}
	}
	return fmt.Sprintf(CONF_TEMP, confFiled)
}
