package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type RedisConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
}

type Conf struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

func loadIni(fileName string, data any) (err error) {
	// 0.参数校验
	// 0.1.data必须是指针，因为要在函数中赋值
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("data must be a pointer")
	}
	// 0.2.data必须是结构体指针，因为配置文件键值对要给结构体
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct pointer")
	}
	// 1.读文件
	b, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	//s := string(b)
	lineSlice := strings.Split(string(b), "\r\n")
	// 2.一行一行读
	var structName string
	for idx, line := range lineSlice {
		// 去掉首尾空格
		line = strings.TrimSpace(line)
		// 2.1.跳过空行
		if len(line) == 0 {
			continue
		}
		// 2.2.如果是#或;表示注释跳过
		if strings.HasPrefix(line, ";") ||
			strings.HasPrefix(line, "#") {
			continue
		}
		// 2.1.如果是[表示节
		if strings.HasPrefix(line, "[") {
			if line[0] != '[' || line[len(line)-1] != ']' {
				return fmt.Errorf("line %d is not a section", idx+1)
			}
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				return fmt.Errorf("line %d is not a section", idx+1)
			}
			// 根据字符串sectionName找对应结构体
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				if sectionName == field.Tag.Get("ini") {
					structName = field.Name
					break
				}
			}
		} else {
			// 2.3.不是[开头就用=分割成键值对
			// 2.3.1.等号分割 左边key 右边value
			if strings.Index(line, "=") == -1 ||
				strings.HasPrefix(line, "=") {
				return fmt.Errorf("line %d is not a key-value", idx+1)
			}
			index := strings.Index(line, "=")
			key := strings.TrimSpace(line[:index])
			val := strings.TrimSpace(line[index+1:])
			// 2.3.2.根据structName 把对应结构体拿出来
			v := reflect.ValueOf(data)
			sVal := v.Elem().FieldByName(structName)
			sType := sVal.Type()

			if sType.Kind() != reflect.Struct {
				return fmt.Errorf("struct %s not found", structName)
			}
			// 2.3.3.根据tag找key
			var fieldName string
			var fieldType reflect.StructField
			for i := 0; i < sVal.NumField(); i++ {
				field := sType.Field(i)
				if field.Tag.Get("ini") == key {
					fieldType = field
					fieldName = field.Name
					break
				}
			}
			if len(fieldName) == 0 {
				return fmt.Errorf("key %s not found", key)
			}
			// 2.3.4.根据fieldName找对应字段
			fieldOvj := sVal.FieldByName(fieldName)
			switch fieldType.Type.Kind() {
			case reflect.String:
				fieldOvj.SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return fmt.Errorf("line %d value %s is not int", idx+1, val)
				}
				fieldOvj.SetInt(valInt)
			case reflect.Bool:
				valBool, err := strconv.ParseBool(val)
				if err != nil {
					return fmt.Errorf("line %d value %s is not bool", idx+1, val)
				}
				fieldOvj.SetBool(valBool)
			case reflect.Float32, reflect.Float64:
				valFloat, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Errorf("line %d value %s is not float", idx+1, val)
				}
				fieldOvj.SetFloat(valFloat)
			}
		}
	}
	return nil
}

func main() {
	var c Conf
	err := loadIni("./conf.ini", &c)
	if err != nil {
		fmt.Println("load ini failed err:", err)
		return
	}
	fmt.Printf("%#v\n", c)
}
