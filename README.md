# Service to manage configuration data: service-config-data

[![GitHub release](https://img.shields.io/github/release/OlegGorj/service-config-data.svg)](https://github.com/OlegGorj/service-config-data/releases)
[![GitHub issues](https://img.shields.io/github/issues/OlegGorj/service-config-data.svg)](https://github.com/OlegGorj/service-config-data/issues)
[![GitHub commits](https://img.shields.io/github/commits-since/OlegGorj/service-config-data/0.0.1.svg)](https://github.com/OlegGorj/service-config-data)

[![Build Status](https://travis-ci.org/OlegGorj/service-config-data.svg?branch=master)](https://travis-ci.org/OlegGorj/service-config-data)
[![GitHub Issues](https://img.shields.io/github/issues/OlegGorJ/service-config-data.svg)](https://github.com/OlegGorJ/service-config-data/issues)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/OlegGorJ/service-config-data.svg)](http://isitmaintained.com/project/OlegGorJ/service-config-data "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/OlegGorJ/service-config-data.svg)](http://isitmaintained.com/project/OlegGorJ/service-config-data "Percentage of issues still open")
[![Docker Stars](https://img.shields.io/docker/stars/oleggorj/service-config-data.svg)](https://hub.docker.com/r/oleggorj/service-config-data/)
[![Docker Pulls](https://img.shields.io/docker/pulls/oleggorj/service-config-data.svg)](https://hub.docker.com/r/oleggorj/service-config-data/)
[![ImageLayers](https://images.microbadger.com/badges/image/oleggorj/service-config-data.svg)](https://microbadger.com/#/images/oleggorj/service-config-data)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1818748c6ba745ce97bb43ab6dbbfd2c)](https://www.codacy.com/app/OlegGorj/service-config-data?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=OlegGorj/service-config-data&amp;utm_campaign=Badge_Grade)


The `service-config-data` service designed to provide configuration management capabilities across all components. This common service can be used by any application, component or module within permitted boundaries. In it's present implementation, the service allows retrieval of values based on provided key and specified environment (`dev`, `stg`, `prod`, etc.).

---

## Directory structure

```
├── service-config-data
    ├── Dockerfile
    ├── Makefile
    ├── README.md
    ├── config-data-util
    │   ├── config_structs.go
    │   ├── data_utility.go
    │   ├── environment
    │   │   └── environment.go
    │   ├── kernel
    │   │   ├── kernel.go
    │   │   └── kernel_test.go
    │   └── user
    │       ├── user.go
    │       └── user_test.go
    ├── credentials_template.json
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
    ├── service.go
    ├── tmp.yaml
    ├── vars.mk
    └── yaml
        ├── deployment.yaml
        ├── secretBluemixDefault.yaml
        ├── secretBluemixDefaultInternational.yaml
        ├── secretBluemixDefaultRegional.yaml
        ├── secretPullImage.yaml
        └── service.yaml
```
---
## TODO: THIS NEEDS TO GET UPDATED
## How to deploy service

### Prerequisites:

The following environment variables must be set as part of `vars.mk` in order for service to compile and run properly:

`vars.mk` file

```
REPO?=github.com/OlegGorJ/config-data
GITACCOUNT?=<your git account>
APITOKEN?=<your API Git token>
CONFIGFILE?=services
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