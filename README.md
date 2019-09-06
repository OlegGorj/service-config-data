# Service to manage configuration data (service-config-data)

[![GitHub release](https://img.shields.io/github/release/OlegGorj/service-config-data.svg)](https://github.com/OlegGorj/service-config-data/releases)
[![GitHub issues](https://img.shields.io/github/issues/OlegGorj/service-config-data.svg)](https://github.com/OlegGorj/service-config-data/issues)
[![GitHub commits](https://img.shields.io/github/commits-since/OlegGorj/service-config-data/0.0.1.svg)](https://github.com/OlegGorj/service-config-data)

[![Build Status](https://travis-ci.org/OlegGorj/service-config-data.svg?branch=master)](https://travis-ci.org/OlegGorj/service-config-data)
[![GitHub Issues](https://img.shields.io/github/issues/OlegGorJ/service-config-data.svg)](https://github.com/OlegGorJ/service-config-data/issues)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/OlegGorJ/service-config-data.svg)](http://isitmaintained.com/project/OlegGorJ/service-config-data "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/OlegGorJ/service-config-data.svg)](http://isitmaintained.com/project/OlegGorJ/service-config-data "Percentage of issues still open")
[![Docker Stars](https://img.shields.io/docker/stars/oleggorj/service-config-data.svg)](https://hub.docker.com/r/oleggorj/service-config-data/)
[![Docker Pulls](https://img.shields.io/docker/pulls/oleggorj/service-config-data.svg)](https://hub.docker.com/r/oleggorj/service-config-data/)


[![Docker Build Status]()](https://hub.docker.com/r/oleggorj/service-config-data/builds/)



[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1818748c6ba745ce97bb43ab6dbbfd2c)](https://www.codacy.com/app/OlegGorj/service-config-data?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=OlegGorj/service-config-data&amp;utm_campaign=Badge_Grade)


The `service-config-data` service designed to provide configuration management capabilities across all components. This common service can be used by any application, component or module within permitted boundaries. In it's present implementation, the service allows retrieval of values based on provided key and specified environment (`dev`, `stg`, `prod`, etc.).

---

## Directory structure

```
── service-config-data
    ├── JenkinsPod.yaml
    ├── Jenkinsfile
    ├── Makefile
    ├── OWNERS
    ├── OWNERS_ALIASES
    ├── README.md
    ├── charts
    │   ├── preview
    │   │   ├── Chart.yaml
    │   │   ├── Makefile
    │   │   ├── requirements.yaml
    │   │   └── values.yaml
    │   └── service-config
    │       ├── Chart.yaml
    │       ├── requirements.yaml
    │       ├── templates
    │       │   ├── deployment.yaml
    │       │   ├── image-pull-secret.yaml
    │       │   └── service.yaml
    │       └── values-template.yaml
    ├── config-data-util
    │   ├── config_structs.go
    │   ├── data_utility.go
    │   ├── environment
    │   │   └── environment.go
    │   ├── kernel
    │   │   ├── kernel.go
    │   │   └── kernel_test.go
    │   ├── memfilesystem
    │   │   └── memfilesystem.go
    │   └── user
    │       ├── user.go
    │       └── user_test.go
    ├── gitutil
    │   ├── git_utility.go
    │   ├── git_utility_test.go
    │   └── test_data
    │       ├── kernels
    │       │   └── spark-kernel-12CPU-24GB
    │       ├── sandbox_whitelist.json
    │       ├── test_object_storage.json
    │       ├── test_users1.json
    │       ├── test_users2.json
    │       └── user.json
    ├── glide.yaml
    ├── handlers
    │   ├── confEnvHandler.go
    │   ├── gitHandler.go
    │   ├── kernelHandler.go
    │   ├── userHandler.go
    │   └── userHandler_test.go
    ├── helpers
    │   └── helpers.go
    ├── service-common-lib
    │   └── common
    │       └── config
    │           └── config.go
    ├── service-config-data
    ├── service.go
    ├── skaffold.yaml
    ├── vars-gcp.mk
    ├── vars.mk
    └── watch.sh
```
---

## How to deploy Config Service

### Prerequisites:

Setup GOPATH:

```
echo $GOPATH
export GOPATH=$GOPATH:$PWD
echo $GOPATH
```

Install common package:

```
go get -u github.com/oleggorj/service-common-lib
```

Install `service-config-data` package:

```
go get -u github.com/oleggorj/service-config-data
```

Build `service-config-data` service binaries:

```
cd src/github.com/oleggorj/service-config-data
make build
```

### Deployment of GCP:

The following environment variables must be set as part of `vars-gcp.mk` in order for service to compile and run properly:

`vars-gcp.mk` file

```
REPO?=github.com/OlegGorJ/config-data
GITACCOUNT?=
APITOKEN?=
CONFIGFILE?=services
REGISTRY?=oleggorj
```

Initialize GCP environment:

```
gcloud auth login
```

List available kubernetes clusters:

```
gcloud container clusters list
```

Example of output:

```
NAME                      LOCATION       MASTER_VERSION  MASTER_IP      MACHINE_TYPE   NODE_VERSION   NUM_NODES  STATUS
cluster-services-sandbox  us-central1-a  1.13.7-gke.19   35.184.197.12  n1-standard-1  1.13.7-gke.19  3          RUNNING
```

Than, create `tiller` service account and cluster buinding:

```
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
```

And, finally, initialize Helm:

```
helm init --service-account tiller --upgrade
```

Last step to build service binaries, create image, push image to docker repo and deploy Helm chart:

```
cd service-config-data
make deploy
```


## Manual step-by-step deployment

### 0. Build service runtime from Go code

```
make build
```

### 1. Build image and push to registry

```
make push
```

### 2. Deploy service

```
make deploy
```

### 3. Clean up

```
make deployclean
```

---
## How to use it.


### 1. `config-data` repository

Configuration file for Services `config-data/service-config-data.json`, branch `sandbox`, looks like

```
{
  "app" : "service-config-data",
  "app-type" : "backend",
  "env" : "sandbox"
}
```

### 2. Call `service-config-data` service

Set environment SHELL variables:

```
export ENVIRONMENT=sandbox
export APP=services
export KEY=app-type
```

Call the service:

```
curl http://<IP of service-config-data service>:8000/api/v1/$APP/$ENVIRONMENT/$KEY
```

Example of calling the Config service from *inside* the k8s cluster:

```
 curl http://service-config-data.default:8000/api/v1/services/sandbox/app
```
You should get back:

```
service-config-data
```
---

# Deploy service on K8S cluster

Make sure values file `./service-config/values.yaml` is populated based on template `./service-config/values-template.yaml`

Deploy service by executing the following

```
cd service-config-data/charts
helm upgrade --install config-service --values ./service-config/values.yaml --namespace default  ./service-config/
```


Clean up:

```
helm del --purge config-service
```


---
### Note for DevOps

*This section needs updates*

_Note:_ internal IPs

Jenkins build job for [Master branch](http://10.0.0.11:8080/job/service-config_master/)

Jenkins build job for [Pull requests](http://10.0.0.11:8080/job/service-config_pr/)

---
