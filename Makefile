COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOOS?=linux
GOARCH?=amd64

include vars-gcp.mk

APP?=service-config-data
APIVER?="v2"

RELEASE?=1.3
IMAGE?=${REGISTRY}/${APP}:${RELEASE}

PORT?=8000
NODE_PORT?=30083

ENV?=SANDBOX

K8S_CHART?=service-config
K8S_NAMESPACE?=default
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
		docker stop ${APP} || true && docker rm ${APP} || true
		docker run --name ${APP} -p ${PORT}:${PORT} --rm \
			-e "PORT=${PORT}" \
			$(IMAGE)

push: container
		docker push $(IMAGE)

container: build
		# generate dockerfile from template
		for t in $(shell find ./src/github.com/oleggorj/service-common-lib/docker/ -type f -name "Dockerfile.goservice.template"); do \
					cat $$t | \
						sed -E "s/{{ .PORT }}/$(PORT)/g" | \
						sed -E "s/{{ .ServiceName }}/$(APP)/g"; \
		done > ./Dockerfile
		docker build -t $(IMAGE) .
		rm Dockerfile
		rm -f ${APP}

deployclean:
		-helm del --purge ${K8S_CHART}

deploy: push
		echo ""
		echo "*** did you run 'make push'? ***"
		echo ""
		for t in $(shell find ./charts/${K8S_CHART} -type f -name "values-template.yaml"); do \
					cat $$t | \
						sed -E "s/{{ .ServiceName }}/$(APP)/g"; \
		done > ./charts/${K8S_CHART}/values.yaml
		cd charts
		helm upgrade --install ${K8S_CHART} --values ./charts/${K8S_CHART}/values.yaml --namespace ${K8S_NAMESPACE}  ./charts/${K8S_CHART}/
		rm ./charts/${K8S_CHART}/values.yaml

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
		curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
		glide install
