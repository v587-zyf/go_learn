package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level LogLv

	filePath string // 日志保存路径
	fileName string // 日志文件名

	fileObj    *os.File
	errFileObj *os.File

	maxFileSize int64
}

func NewFileLogger(lv, fp, fn string, mfs int64) *FileLogger {
	f := &FileLogger{
		Level:       str2Lv(lv),
		filePath:    fp,
		fileName:    fn,
		maxFileSize: mfs,
	}

	err := f.initFile()
	if err != nil {
		panic(err)
	}

	return f
}

func (f *FileLogger) initFile() error {
	fName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open log file err:", err)
		return err
	}

	fName = path.Join(f.filePath, "err_"+f.fileName)
	errFileObj, err := os.OpenFile(fName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open err log file err:", err)
		return err
	}

	f.fileObj = fileObj
	f.errFileObj = errFileObj

	return nil
}

func (f *FileLogger) enable(logLv LogLv) bool {
	return f.Level <= logLv
}

func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info err:", err)
		return false
	}

	return fileInfo.Size() >= f.maxFileSize
}

func (f *FileLogger) logPrint(lv LogLv, format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	funcName, fileName, line := getInfo(3)
	if f.checkSize(f.fileObj) {
		newFile, err := f.SplitFile(f.fileObj)
		if err != nil {
			fmt.Println("split log file err:", err)
			return
		}
		f.fileObj = newFile
	}
	fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n",
		getTime(), lv2Str(lv),
		fileName, funcName, line, str)
	if lv >= ERROR {
		if f.checkSize(f.errFileObj) {
			newFile, err := f.SplitFile(f.errFileObj)
			if err != nil {
				fmt.Println("split err log file err:", err)
				return
			}
			f.errFileObj = newFile
		}
		fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n",
			getTime(), lv2Str(lv),
			fileName, funcName, line, str)
	}
}

func (f *FileLogger) SplitFile(file *os.File) (*os.File, error) {
	// 切割文件
	// 2.备份 rename
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info err:", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())
	newFileName := fmt.Sprintf("%s.bak%s", logName, nowStr)

	// 1.关闭当前文件
	file.Close()

	os.Rename(logName, newFileName)
	// 3.打开新文件
	fileObj, err := os.OpenFile(logName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open logName file err:", err)
		return nil, err
	}
	return fileObj, nil
}

func (f *FileLogger) Debug(str string, a ...any) {
	lv := DEBUG
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Trace(str string, a ...any) {
	lv := TRACE
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Info(str string, a ...any) {
	lv := INFO
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Warn(str string, a ...any) {
	lv := WARN
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Error(str string, a ...any) {
	lv := ERROR
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Fatal(str string, a ...any) {
	lv := FATAL
	if !f.enable(lv) {
		return
	}
	f.logPrint(lv, str, a...)
}

func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
