SET TAG=dev
::SET REPO=192.168.61.190:5000
SET REPO=47.236.235.165:5000

cd login
make build tag=dev repo=47.236.235.165:5000 df=Dockerfile && make push tag=dev repo=47.236.235.165:5000 df=Dockerfile
cd ../

rem cd gate
rem make build tag=dev repo=47.236.235.165:5000 df=Dockerfile && make push tag=dev repo=47.236.235.165:5000 df=Dockerfile
rem cd ../

rem center game gate login register
