@echo off

::set PROJECT_ROOT_DIR=D:\code\golearn
set PROJECT_ROOT_DIR=D:\code\go\mokoko\thunder\comm
::生成objs.go文件的路径
::set savePath=%PROJECT_ROOT_DIR%\config
set savePath=%PROJECT_ROOT_DIR%\t_tdb\
::目标excel文件路径
::set confPath=%PROJECT_ROOT_DIR%
set confPath=%PROJECT_ROOT_DIR%\t_config
::所有的字段类型
set allType=int,IntSlice,IntSlice2,IntMap,string,StringSlice,StringSlice2,float64,HmsTime,HmsTimes,bool,FloatSlice
::package
set package=t_tdb

..\bin\genStruct.exe -savePath=%savePath% -readPath=%confPath%\excel -allType=%allType% -package=%package%

::pause
::TIMEOUT /T 10