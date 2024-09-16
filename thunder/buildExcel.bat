@echo off

echo update online start...

set UP=true
set IP="47.236.235.165"
set PORT=62222
set USER="root"
set PASS="3FlKdGkPycJdfgpn"
set CMD=true
set CMD_INFO="cd /data/go/thunder/config&&rm tdb.dat"
set LP="./excel"
set RP="/data/go/thunder/config"
set F=".git;.gitignore"
set TOOL_DIR="D:\code\go\mokoko\bin\upload.exe"

cd ./comm/t_config
%TOOL_DIR% -t=s -up=%UP% -ip=%IP% -port=%PORT% -user=%USER% -pass=%PASS% -c=%CMD% -ci=%CMD_INFO% -lp=%LP% -rp=%RP% -f=%F%
rem cd  ..\
echo update online end...

::pause
::TIMEOUT /T 3