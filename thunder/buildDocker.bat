SET TAG=dev
::SET REPO=192.168.61.190:5000
SET REPO=47.236.235.165:5000

cd center
make build tag=%TAG% repo=%REPO% df=Dockerfile && make push tag=%TAG% repo=%REPO% df=Dockerfile
cd ../

cd game
make build tag=%TAG% repo=%REPO% df=Dockerfile && make push tag=%TAG% repo=%REPO% df=Dockerfile
cd ../

cd gate
make build tag=%TAG% repo=%REPO% df=Dockerfile && make push tag=%TAG% repo=%REPO% df=Dockerfile
cd ../

cd login
make build tag=%TAG% repo=%REPO% df=Dockerfile && make push tag=%TAG% repo=%REPO% df=Dockerfile
cd ../

rem cd register
rem make build tag=%TAG% repo=%REPO% df=Dockerfile && make push tag=%TAG% repo=%REPO% df=Dockerfile
rem cd ../

