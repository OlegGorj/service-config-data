

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
CONTAINER_IMAGE?=registry.ng.bluemix.net/${REGISTRY_NAMESPACE}/${APP}:${RELEASE}
GOOS?=linux
GOARCH?=amd64

#include ./service-common-lib/Makefile

APP?=service-config-data
APIVER?="v2"
REGISTRY?=registry.ng.bluemix.net
REGISTRY_NAMESPACE?=etl-namespace
PORT?=8000
NODE_PORT?=30083
ENV?=SANDBOX
K8S_NAMESPACE?=default
RELEASE?=1.1
REPO?=github.ibm.com/AdvancedAnalyticsCanada/config-data.git

CONTAINER_IMAGE?=${REGISTRY}/${REGISTRY_NAMESPACE}/${APP}:${RELEASE}
NODESELECTOR?=services

clean:
		rm -f ${APP}

build: clean
		echo "GOPATH: " ${GOPATH}
		CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
			-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
			-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
			-o ${APP}

run: container
		docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
		docker run --name ${APP} -p ${PORT}:${PORT} --rm \
			-e "PORT=${PORT}" \
			$(APP):$(RELEASE)

push:
		docker push $(CONTAINER_IMAGE)


container: build
		for t in $(shell find ./service-common-lib/docker/ -type f -name "Dockerfile.goservice.template"); do \
					cat $$t | \
						sed -E "s/{{ .PORT }}/$(PORT)/g" | \
						sed -E "s/{{ .ServiceName }}/$(APP)/g"; \
		done > ./Dockerfile
		docker build -t $(CONTAINER_IMAGE) .

deployclean:
		-helm del --purge config-service

deploy:
		echo ""
		echo "*** did you run 'make push'? ***"
		echo ""
		cd charts
		helm upgrade --install config-service --values ./service-config/values.yaml --namespace default  ./service-config/

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
		curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
		glide install
