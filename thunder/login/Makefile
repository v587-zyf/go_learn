IMAGE_REPO = 127.0.0.1:5000
IMAGE_NAME = login

# rewrite IMAGE_PREFIX
ifneq ($(origin repo), undefined)
IMAGE_REPO = ${repo}
endif


ifneq ($(origin df), undefined)
DOCKERFILE = ${df}
endif

# This is how we want to name the binary output
OUTPUT=../bin/${IMAGE_NAME}

# These are the values we want to pass for Version and BuildTime
GIT_TAG=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`
GIT_HASH=`git rev-parse --short HEAD`
IMAGE_TAG=${tag}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X ${MODULE}/src/conf.GitTag=${GIT_TAG} -X ${MODULE}/src/conf.BuildTime=${BUILD_TIME}"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT} ./main.go && \
#	docker rmi ${IMAGE_NAME}:${IMAGE_TAG} && \
#	docker rmi ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG} && \
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f ./${DOCKERFILE} ../

push:
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG} && \
	docker push ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG} && \
	docker rmi ${IMAGE_NAME}:${IMAGE_TAG} && \
	docker rmi ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG}
