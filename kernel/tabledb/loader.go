package tabledb

import (
	"fmt"
	"kernel/iface"
	"reflect"
)

//	func GetFileInfos() []FileInfo {
//		return fileInfos
//	}
//
//	func Load(basePath string) error {
//		f, err := os.Stat(basePath)
//		if err != nil {
//			return err
//		}
//		isCheck := false
//		InitTableDb(basePath)
//		if f.IsDir() {
//			if err := LoadExcel(basePath, "tableDb.dat"); err != nil {
//				return err
//			}
//			isCheck = true
//		} else {
//			if err := LoadGob(basePath); err != nil {
//				return err
//			}
//			basePath = filepath.Dir(basePath)
//		}
//		//数据整理
//		GetDb().Patch()
//		//数据检查
//		if isCheck {
//			err = GetDb().Check()
//			if err != nil {
//				return err
//			}
//		}
//		//fmt.Println("加载TableDb成功")
//		return nil
//	}
//
//	func Reload() {
//		path := tableDb.TableDbPath
//		if len(path) <= 0 {
//			fmt.Println("配置路径没找到！！！")
//			return
//		}
//		if err := LoadGob(path); err != nil {
//			fmt.Printf("重加载表配置错误：%v", err)
//		}
//		//数据整理
//		GetDb().Patch()
//	}
//
//	func LoadExcel(basePath string, godname string) error {
//		//如果存在dat文件,则先载入.然后对比修改时间,重新加载时间不一致的文件
//		godfile := filepath.Join(basePath, godname)
//		f, err := os.Stat(godfile)
//		if err == nil && !f.IsDir() {
//			err = tableDb.loadGob(godfile)
//			if err != nil {
//				fmt.Println("load gob file error ", err)
//				InitTableDb(godfile)
//			}
//		}
//
//		change, err := tableDb.load(filepath.Join(basePath, "excel"))
//		if err != nil {
//			return err
//		}
//
//		if change {
//			tableDb.Ver = time.Now().Format("060102150405")
//			tableDb.generateGob(godfile)
//		}
//		return nil
//	}
//
//	func LoadGob(filename string) error {
//		err := tableDb.loadGob(filename)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func checkUnique(FileName string, keys map[int64]bool, objV reflect.Value) error {
//		for _, v := range []string{"Id", "Lvl"} {
//			keyFieldV := objV.Elem().FieldByName(v)
//			if keyFieldV.IsValid() {
//				key := keyFieldV.Int()
//				if keys[key] {
//					return fmt.Errorf("表　%s 字段 %s, %d,%v 重复了", FileName, v, key, objV)
//				}
//				keys[key] = true
//				break
//			}
//		}
//		return nil
//
// }
//
//	func arrayLoader(fieldName string) func(*TableDb, []interface{}) error {
//		return func(TableDb *TableDb, objs []interface{}) error {
//			fieldV := reflect.ValueOf(TableDb).Elem().FieldByName(fieldName)
//
//			keys := make(map[int64]bool)
//			if fieldV.Kind() == reflect.Slice {
//				if fieldV.IsNil() || fieldV.Len() > 0 {
//					fieldV.Set(reflect.MakeSlice(fieldV.Type(), 0, len(objs)))
//				}
//				for _, obj := range objs {
//					objV := reflect.ValueOf(obj)
//					fieldV.Set(reflect.Append(fieldV, objV))
//					if err := checkUnique(fieldName, keys, objV); err != nil {
//						return err
//					}
//				}
//			} else if fieldV.Kind() == reflect.Array {
//				for i, obj := range objs {
//					objV := reflect.ValueOf(obj)
//					fieldV.Index(i).Set(objV)
//					if err := checkUnique(fieldName, keys, objV); err != nil {
//						return err
//					}
//				}
//			} else {
//				return fmt.Errorf("field %s is not an array", fieldName)
//			}
//			return nil
//		}
//	}
func MapLoader(fieldName string, keyFieldName string) func(iface.ITableDb, []interface{}) error {
	return func(TableDb iface.ITableDb, objs []interface{}) error {
		fieldV := reflect.ValueOf(TableDb).Elem().FieldByName(fieldName)
		if fieldV.Kind() != reflect.Map {
			return fmt.Errorf("field %s is not a map", fieldName)
		}

		if fieldV.IsNil() || fieldV.Len() > 0 {
			fieldV.Set(reflect.MakeMap(fieldV.Type()))
		}
		for _, obj := range objs {
			objV := reflect.ValueOf(obj)
			keyFieldV := objV.Elem().FieldByName(keyFieldName)
			if !keyFieldV.IsValid() {
				return fmt.Errorf("key field %s wrong filedV:%v, when setting %s\n", keyFieldName, fieldName, keyFieldV)
			}
			if keyFieldV.Kind() != reflect.Int {
				fmt.Printf("key field %s wrong filedV:%v, when setting %s\n", keyFieldName, fieldName, keyFieldV)
				continue
			}
			if fieldV.MapIndex(keyFieldV).IsValid() {
				//return fmt.Errorf(" >%v<. The value of field >%s< in sheet >%s< is duplicate",
				return fmt.Errorf("表 %s 列 %s 值->%v 重复了",
					fieldName, keyFieldName, keyFieldV)
			}
			fieldV.SetMapIndex(keyFieldV, objV)
		}
		return nil
	}
}

//func (t *TableDb) generateGob(FileName string) error {
//
//	//now := time.Now()
//	//defer func() {
//	//	fmt.Println("generateGob use time", time.Since(now).Seconds())
//	//}()
//
//	f, err := os.Create(FileName)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	w := bufio.NewWriter(f)
//	enc := gob.NewEncoder(w)
//	enc.Encode(t)
//	return w.Flush()
//}
//
//func (t *TableDb) loadGob(FileName string) error {
//	//now := time.Now()
//	//defer func() {
//	//	fmt.Println("loadGod use time", time.Since(now).Seconds())
//	//}()
//
//	f, err := os.Open(FileName)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	r := bufio.NewReader(f)
//	dec := gob.NewDecoder(r)
//	return dec.Decode(&t)
//}
//
//func (t *TableDb) load1(baseDir string) (bool, error) {
//	//now := time.Now()
//	//defer func() {
//	//	fmt.Println("load all use time", time.Since(now).Seconds())
//	//}()
//
//	var wg sync.WaitGroup
//	var loadErr error
//	var num int
//	for i := range fileInfos {
//		filename := filepath.Join(baseDir, fileInfos[i].FileName)
//		finfo, err := os.Stat(filename)
//		if err != nil {
//			loadErr = errors.New("filename not found, " + err.Error())
//			break
//		}
//		fmodtime := finfo.ModTime().UnixNano()
//		if t.FileModTime[fileInfos[i].FileName] == fmodtime {
//			//fmt.Printf("文件 %v 未修改!\n", fileInfos[i].FileName)
//			continue
//		}
//		t.FileModTime[fileInfos[i].FileName] = fmodtime
//		num++
//
//		wg.Add(1)
//		go func(index int, filename string) {
//			err := t.loadFile(filename, fileInfos[index].SheetInfos)
//			if err != nil {
//				fmt.Println("加载:", fileInfos[index].FileName, "error:", err)
//				loadErr = err
//			}
//			//fmt.Printf("加载完成: %v\n", fileInfos[index].FileName)
//			wg.Done()
//		}(i, filename)
//	}
//
//	wg.Wait()
//	if loadErr != nil {
//		return false, loadErr
//	}
//	//fmt.Printf("共加载: %v 个文件", num)
//	if num == 0 {
//		return false, nil
//	}
//	//fmt.Println("All TableDb rem cal is ok")
//	return true, loadErr
//}
//
//func (t *TableDb) loadFile(filename string, SheetInfos []SheetInfo) error {
//	xlsFile, err := xlsx.OpenFile(filename)
//	if err != nil {
//		return err
//	}
//	for _, SheetInfo := range SheetInfos {
//		sheet, ok := xlsFile.Sheet[SheetInfo.SheetName]
//		if !ok {
//			return fmt.Errorf("no %s sheet found", SheetInfo.SheetName)
//		}
//		objProptype := reflect.New(reflect.TypeOf(SheetInfo.ObjProptype)).Interface()
//		objs, err := ReadXlsxSheet(sheet, objProptype, 2, 1, nil) // read from 3rd line,2nd row
//		if err != nil {
//			return err
//		}
//		err = SheetInfo.Initer(t, objs)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (t *TableDb) LoadGlobalConf(objs []interface{}) error {
//	tableConfs := make(map[string]*GlobalBaseCfg)
//	for _, obj := range objs {
//		game := obj.(*GlobalBaseCfg)
//		if _, ok := tableConfs[game.Name]; ok {
//			return errors.New(fmt.Sprintf("tableconf key:%d namd:%s 重复了", game.Id, game.Name))
//		}
//		tableConfs[game.Name] = game
//	}
//	return DecodeConfValues(GetConf(), tableConfs)
//}
