@echo off

:: go build -o D:\code\go\mokoko\thunder\bin\genStruct.exe .
:: go build -o D:\code\go\mokoko\bar\bin\genStruct.exe .
:: go build -o D:\code\go\mokoko\bin\genStruct.exe .

::set PROJECT_ROOT_DIR=D:\code\golearn
set PROJECT_ROOT_DIR=D:\code\mokoko\thunder
::生成objs.go文件的路径
::set savePath=%PROJECT_ROOT_DIR%\config
set savePath=%PROJECT_ROOT_DIR%\config\config
::目标excel文件路径
::set confPath=%PROJECT_ROOT_DIR%
set confPath=%PROJECT_ROOT_DIR%
::所有的字段类型
set allType=int,IntSlice,IntSlice2,IntSlice3,IntMap,string,StringSlice,StringSlice2,float64,ItemInfo,ItemInfos,PropInfo,PropInfos,HmsTime,HmsTimes,bool,FloatSlice
::package
set package=tabledb

.\genStruct.exe -savePath=%savePath% -readPath=%confPath%\excel -allType=%allType% -package=%package%

TIMEOUT /T 10