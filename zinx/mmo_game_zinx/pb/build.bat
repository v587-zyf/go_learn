@echo off

protoc --go_out=. *.proto
rem protoc --gofast_out=. *.proto