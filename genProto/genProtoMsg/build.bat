@echo off

set BUILD_DIR=D:/code/mokoko/thunder/bin

go build -o %BUILD_DIR%%/genProto.exe .

echo build to "%BUILD_DIR%"

TIMEOUT /T 3

:: go build -o D:/code/go/gc/examples/protobuf/genProto.exe .