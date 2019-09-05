COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOOS?=linux
GOARCH?=amd64

include vars-gcp.mk

APP?=service-config-data
APIVER?="v2"
#REGISTRY?=registry.ng.bluemix.net
#REGISTRY_NAMESPACE?=etl-namespace

RELEASE?=1.0
IMAGE?=${REGISTRY}/${APP}:${RELEASE}

PORT?=8000
NODE_PORT?=30083

ENV?=SANDBOX

K8S_CHART?=config-service
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
		docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
		docker run --name ${APP} -p ${PORT}:${PORT} --rm \
			-e "PORT=${PORT}" \
			$(APP):$(RELEASE)

push: container
		docker push $(IMAGE)

container: build
	# generate dockerfile from template
		for t in $(shell find ../service-common-lib/docker/ -type f -name "Dockerfile.goservice.template"); do \
					cat $$t | \
						sed -E "s/{{ .PORT }}/$(PORT)/g" | \
						sed -E "s/{{ .ServiceName }}/$(APP)/g"; \
		done > ./Dockerfile
		docker build -t $(IMAGE) .

deployclean:
		-helm del --purge ${K8S_CHART}

deploy:
		echo ""
		echo "*** did you run 'make push'? ***"
		echo ""
		# TODO generate values.yaml file from template
		cd charts
		helm upgrade --install ${K8S_CHART} --values ./charts/service-config/values.yaml --namespace ${K8S_NAMESPACE}  ./charts/service-config/

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
		curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
		glide install
