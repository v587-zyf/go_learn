@echo off

set ROOT_DIR=D:/code/go/mokoko/thunder

set TOOL_DIR=D:/code/go/mokoko/bin
set PB_DIR=%ROOT_DIR%/comm/t_proto

set DIR_SOURCE=%PB_DIR%/source
set DIR_OUT=%PB_DIR%/out

set GEN_DIR=%TOOL_DIR%/genProto.exe

set PROTO_PATH=C:/Users/Administrator/go/pkg/mod/github.com/google/gnostic@v0.5.7-v3refs/third_party/

cd %DIR_SOURCE%
for /f "delims=\" %%f in ('dir /b "*.*"') do (
    if /i not "%%f"==".gitignore" if /i not "%%f"=="README.md" (
        cd %%f

        echo genProtoMsg %DIR_SOURCE%/%%f/ start...
        %GEN_DIR% -s %DIR_SOURCE%/%%f/ -o %DIR_OUT%/%%f/
        echo genProtoMsg %DIR_SOURCE%/%%f/ end...

        echo genProto %DIR_SOURCE%/%%f/ start...
        protoc --proto_path=%PROTO_PATH% -I=%DIR_SOURCE%/%%f/ --gofast_out=%DIR_OUT%/%%f/ %DIR_SOURCE%/%%f/*.proto
        protoc --proto_path=%PROTO_PATH% -I=%DIR_SOURCE%/%%f/ --gofast_out=plugins=grpc:%DIR_OUT%/%%f/ %DIR_SOURCE%/%%f/*.proto

        ::protoc -I=%DIR_SOURCE%/%%f/ --go_out=%DIR_OUT%/%%f/ %DIR_SOURCE%/%%f/*.proto
        echo genProto %DIR_SOURCE%/%%f/ end...
        echo.
        echo.

        cd ../
    )
)

rem pause
rem TIMEOUT /T 3