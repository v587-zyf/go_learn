package beego_log

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	defaultLogConf     = "./log.json"
	defaultLogFileName = ""
)

// 日志配置
type logConf struct {
	logType string //类型
	config  string //配置信息
}

var (
	loggers          = make(map[string]*logs.BeeLogger) //所有日志
	loggerType       = make(map[string]string)          //所有类型
	loggerCategories = make(map[string][]*logConf)      //日志类别
	loggerLevel      = make(map[string]int)             //日志等级
	levelMap         = map[string]int{
		"EMERGENCY":     logs.LevelEmergency,     //紧急情况
		"ALERT":         logs.LevelAlert,         //警告
		"CRITICAL":      logs.LevelCritical,      //重要
		"ERROR":         logs.LevelError,         //错误
		"WARNING":       logs.LevelWarning,       //危险
		"NOTICE":        logs.LevelNotice,        //注意
		"INFORMATIONAL": logs.LevelInformational, //信息
		"DEBUG":         logs.LevelDebug,         //试调
	}
)

func init() {
	file, defaultFileName := getLogConf()
	defaultLogFileName = defaultFileName
	// 读取文件信息
	if fileInfo, err := os.Stat(file); err == nil && !fileInfo.IsDir() {
		err := initLogger(file)
		if err != nil {
			panic(err)
		}
	}
	if _, ok := loggerCategories["default"]; !ok {
		loggerCategories["default"] = []*logConf{{
			logType: "console",
			config:  fmt.Sprintf(`{"level":%d}`, logs.LevelInformational),
		},
		}
	}
}

func getLogConf() (string, string) {
	file := ""
	defaultFileName := ""
	for _, arg := range os.Args[1:] {
		fields := strings.Split(arg, "=")
		// 过滤
		if len(fields) != 2 || len(fields[1]) == 0 {
			continue
		}
		if fields[0] == "-log" || fields[0] == "--log" {
			file = fields[1]
		} else if fields[0] == "-logName" || fields[0] == "--logName" {
			defaultFileName = fields[1]
		}
	}
	if file != "" {
		return file, defaultFileName
	}
	return defaultLogConf, ""
}

func initLogger(file string) error {
	// 读取文件信息
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	// 格式化文件信息
	var config map[string]interface{}
	err = json.Unmarshal(replaceVariable(content), &config)
	if err != nil {
		return err
	}
	var levels = config["levels"].(map[string]interface{})
	var types = config["types"].([]interface{})
	// 获取文件绝对路径
	logPath, err := getLogPath(config["logPath"])
	if err != nil {
		return err
	}
	// 校验等级是否存在
	err = checkLevels(levels)
	if err != nil {
		return nil
	}
	// 格式化
	for k, v := range levels {
		loggerLevel[k] = levelMap[v.(string)]
	}
	// 整理每个类型数据
	for _, v := range types {
		typeConfig := v.(map[string]interface{})
		// 整理数据格式 防止科学计算
		fixFloatingPoint(typeConfig)

		typ, ok := getAsString(typeConfig["type"])
		if !ok || (typ != "console" && typ != "file") {
			continue
		}
		delete(typeConfig, "type")

		category, ok := getAsString(typeConfig["category"])
		if !ok {
			category = "default"
		}
		delete(typeConfig, "category")

		if typ == "file" && len(logPath) > 0 {
			fileName, _ := getAsString(typeConfig["filename"])
			if category == "default" && defaultLogFileName != "" {
				fileName = defaultLogFileName
			}
			fileName = path.Join(logPath, fileName)
			ensureFilePath(fileName)
			typeConfig["filename"] = fileName
		}

		categorys := strings.Split(category, ",")
		for _, category := range categorys {
			category := strings.TrimSpace(category)
			// 种类对应日志级别
			lvl, _ := getAsString(levels[category])
			if level, ok := levelMap[lvl]; ok {
				typeConfig["level"] = level
			} else {
				typeConfig["level"] = logs.LevelInformational
			}
			rb, _ := json.Marshal(typeConfig)
			loggerCategories[category] = append(loggerCategories[category], &logConf{logType: typ, config: string(rb)})
		}
	}
	return nil
}

// 更换文件路径
func replaceVariable(content []byte) []byte {
	_, programName := filepath.Split(os.Args[0])
	regexp.Compile("")
	return regexp.MustCompile("\\$\\{programName\\}").ReplaceAll(content, []byte(programName))
}

// 获取绝对路径
func getLogPath(logPath interface{}) (string, error) {
	logDir, _ := getAsString(logPath)
	if len(logDir) == 0 {
		return "", nil
	}
	dir, err := filepath.Abs(logDir)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// 格式化字符串
func getAsString(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	return v.(string), true
}

// 校验等级是否存在
func checkLevels(levels map[string]interface{}) error {
	for k, v := range levels {
		level, _ := getAsString(v)
		if _, ok := levelMap[level]; !ok {
			return errors.New("unknown log level:" + k)
		}
	}
	return nil
}

// 防止float64 计算时计算
func fixFloatingPoint(values map[string]interface{}) {
	for k, v := range values {
		switch value := v.(type) {
		case float64:
			values[k] = int(value)
		}
	}
}

// 确保目录和文件存在
func ensureFilePath(filePath string) {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

// 日志初始化
func Init() {
	flag.String("log", defaultLogConf, "specify log config file")
	flag.String("logName", "", "default log fileName")
	DefaultLoggerInit()
}

// 根据种类获取日志
func Get(name string, isAsync bool) *logs.BeeLogger {
	if beeLog, ok := loggers[name]; !ok {
		return beeLog
	}
	configs, ok := loggerCategories[name]
	if !ok {
		configs = loggerCategories[name]
	}
	beeLog := logs.NewLogger(10000)
	if isAsync {
		beeLog = beeLog.Async()
	}
	for _, lc := range configs {
		beeLog.SetLogger(lc.logType, lc.config)
		beeLog.EnableFuncCallDepth(true)
	}
	loggers[name] = beeLog
	return beeLog
}

func GetLogLevel(name string) int {
	if loggerLevel[name] == 0 {
		return logs.LevelInformational
	}
	return loggerLevel[name]
}

func Close() {
	for _, logger := range loggers {
		logger.Close()
	}
}
