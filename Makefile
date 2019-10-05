COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOOS?=linux
GOARCH?=amd64

include vars-gcp.mk

APP?=service-config-data
APIVER?=v2
RELEASE?=1.5
IMAGE?=${DOCKER_ORG}/${APP}:${RELEASE}

ENV?=DEV

K8S_CHART?=service-config
K8S_NAMESPACE?=dev
NODESELECTOR?=services

helm:
		kubectl create serviceaccount --namespace kube-system tiller
		kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
		helm init --service-account tiller --upgrade

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
		helm del --purge "${K8S_CHART}-${K8S_NAMESPACE}"

deploy:
		for t in $(shell find ./charts/${K8S_CHART} -type f -name "values-template.yaml"); do \
					cat $$t | \
						sed -E "s/{{ .ServiceName }}/$(APP)/g" | \
						sed -E "s/{{ .Release }}/$(RELEASE)/g" | \
						sed -E "s/{{ .Env }}/$(ENV)/g" | \
						sed -E "s/{{ .Kube_namespace }}/$(K8S_NAMESPACE)/g" | \
						sed -E "s/{{ .ApiVer }}/$(APIVER)/g" | \
						sed -E "s/{{ .LBPort }}/$(LB_EXTERNAL_PORT)/g" | \
						sed -E "s/{{ .ContainerPort }}/$(PORT)/g" | \
						sed -E "s/{{ .DockerOrg }}/$(DOCKER_ORG)/g"; \
		done > ./charts/${K8S_CHART}/values.yaml
		helm install --name "${K8S_CHART}-${K8S_NAMESPACE}" --values ./charts/${K8S_CHART}/values.yaml --namespace ${K8S_NAMESPACE}  ./charts/${K8S_CHART}/
		echo "Cleaning up temp files.." && rm ./charts/${K8S_CHART}/values.yaml
		kubectl get services --all-namespaces | grep ${APP}
		./scripts/githook.sh ${APP} ${LB_EXTERNAL_PORT} webhook_git ${GITUSER} ${GITREPO} ${K8S_NAMESPACE}


.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
		curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
		glide install
