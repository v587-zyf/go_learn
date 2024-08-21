@echo off

set ROOT_DIR=D:/main/golearn

set TOOL_DIR=%ROOT_DIR%/tools/genProto/
set PB_DIR=%ROOT_DIR%/protobuf/

set DIR_SOURCE=%TOOL_DIR%source/
set DIR_OUT=%PB_DIR%out/

set GEN_DIR=%TOOL_DIR%/genProto.exe

cd %DIR_SOURCE%
for /f "delims=\" %%f in ('dir /b "*.*"') do (
    cd %%f

    echo ******************
    %GEN_DIR% -s %DIR_SOURCE%%%f/ -o %DIR_OUT%%%f/
    echo generate %%f const and idMap end

    echo ------------------
    protoc -I=%DIR_SOURCE%%%f/ --gofast_out=%DIR_OUT%%%f/ %DIR_SOURCE%%%f/*.proto
    echo generate %%f proto end

    cd ../
)

rem pause
TIMEOUT /T 10

rem protoc -I=D:/code/golearn/demo/grpc/proto/ --gofast_out=./ ./*.proto