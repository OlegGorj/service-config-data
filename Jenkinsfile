#!/usr/bin/env groovy

import groovy.json.JsonOutput
import hudson.model.*
import java.util.regex.Pattern

def K8S_CLUSTER_NAME
def CONFIG_SERVICE = "service-config-data.default"
def ENV = "sandbox"
def ORG = "AdvancedAnalyticsCanada"
def KUBECTL_VERSION = "v1.13.0"
def HELM_VERSION = "v2.1.3"
def DOCKER_VERSION = "1.12.6"

def secret_access_key = ""
def GIT_AUTH_ID = ""
def GIT_IBM_AUTH_TOKEN = ""
def BM_APIKEY = ""
def BM_API_ENDPOINT = "https://api.us-east.bluemix.net"
def WORK_DIR = ""

def SLACK_TOKEN
def SLACK_URL
def SLACK_CHANNEL
def smoke_test_status
def perf_test_status
def integr_test_status
def K8S_NAMESPACE

def envs = ['sandbox','dev','staging','prod']


properties([
    parameters([
        choice(choices: envs, description: 'Please select an environment', name: 'Environments'),
    ]),
    pipelineTriggers([

    ])
])

podTemplate(
    label: 'build-pod',
    containers: [
        containerTemplate(name: 'jnlp', image: 'jenkins/jnlp-slave:3.10-1-alpine', args: '${computer.jnlpmac} ${computer.name}'),
        containerTemplate(
          name: 'golang', image: 'golang:1.12-alpine3.9', ttyEnabled: true, command: 'cat',
          envVars: [
            envVar(key: 'KUBECONFIG', value: '/home/jenkins/.bluemix/plugins/container-service/clusters/k8s-husky-sandbox/kube-config-tor01-k8s-husky-sandbox.yml'),
            envVar(key: 'ORG', value: 'AdvancedAnalyticsCanada'),
            envVar(key: 'APP_NAME', value: 'service-config-data'),
            envVar(key: 'PREVIEW_VERSION', value: '0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER'),
            envVar(key: 'PREVIEW_NAMESPACE', value: '$APP_NAME-$BRANCH_NAME'),
            envVar(key: 'HELM_RELEASE', value: '$PREVIEW_NAMESPACE')
          ]
        ),
        /*
        containerTemplate(
          name: 'maven-custom', image: 'jenkinsxio/builder-maven:0.1.235', ttyEnabled: true, command: 'cat',
          envVars: [
            envVar(key: 'KUBECONFIG', value: '/home/jenkins/.bluemix/plugins/container-service/clusters/k8s-husky-sandbox/kube-config-tor01-k8s-husky-sandbox.yml'),
            envVar(key: 'ORG', value: 'AdvancedAnalyticsCanada'),
            envVar(key: 'APP_NAME', value: 'service-config-data'),
            envVar(key: 'PREVIEW_VERSION', value: '0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER'),
            envVar(key: 'PREVIEW_NAMESPACE', value: '$APP_NAME-$BRANCH_NAME'),
            envVar(key: 'HELM_RELEASE', value: '$PREVIEW_NAMESPACE')
          ]
        ),*/
    ],
    volumes: [
        hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
        secretVolume(mountPath: '/secrets/auth', secretName: 'auth-pipeline-secrets')
    ],
    imagePullSecrets: [ 'image-pull-secret-ibm-cloud' ],
    namespace: "jx",
    nodeSelector: 'worker: true'
) {

    node('build-pod') {
        container(name: 'jnlp') {
          echo "Running build ${env.BUILD_ID}, branch $env.BRANCH_NAME on ${env.JENKINS_URL}"
        }

        container(name: 'golang') {

        try {
        stage('Install packages') {
            sh '''
              #curl -Lo /tmp/docker.tgz https://get.docker.com/builds/Linux/x86_64/docker-1.12.6.tgz
              #mkdir /tmp/docker
              #tar -xf /tmp/docker.tgz -C /tmp/docker
              #ls -l /tmp/docker/docker/
              #mv /tmp/docker/docker/docker* /usr/local/bin/
              #apt-get purge docker lxc-docker docker-engine docker.io
              #apt-get update
              #apt-get install -y --force-yes apt-transport-https software-properties-common
              #curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add
              #add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
              #apt-get update || true
              #apt-get install docker-ce
              apk update --no-cache && apk upgrade --no-cache
              apk add openrc --no-cache
              apk add git make bash curl jq g++ ca-certificates docker
              apk add device-mapper
              rc-update add docker boot
              git clone https://github.com/jpetazzo/dind.git
              cp dind/wrapdocker /usr/local/bin/wrapdocker
              chmod +x /usr/local/bin/wrapdocker
              /usr/local/bin/wrapdocker &
              service docker start || true
              service docker status || true
            '''
            }

            stage('Init Stage') {
              GIT_AUTH_ID = sh(script: "cat /secrets/auth/GIT_AUTH_ID", returnStdout: true)
              GIT_IBM_AUTH_TOKEN = sh(script: "cat /secrets/auth/GIT_IBM_AUTH_TOKEN", returnStdout: true)
              //sh 'export GITHUB_TOKEN=`cat /secrets/auth/GIT_IBM_AUTH_TOKEN`'
              //sh 'export GIT_IBM_AUTH_TOKEN=`cat /secrets/auth/GIT_IBM_AUTH_TOKEN`'
              //BM_APIKEY = sh(script: "cat /secrets/auth/BM_APIKEY", returnStdout: true).toString().trim()
              //sh 'export IBMCLOUD_API_KEY=`cat /secrets/auth/BM_APIKEY`'
              K8S_CLUSTER_NAME = sh(script: "curl http://${CONFIG_SERVICE}:8000/api/v1/k8s-cluster/${ENV}/k8s-cluster-name", returnStdout: true).toString().trim()
            }

            stage('Setup/Checkout SCM') {
              script {
                sh '''
                  export GITHUB_TOKEN=`cat /secrets/auth/GIT_IBM_AUTH_TOKEN`
                  git config --global http.extraheader "PRIVATE-TOKEN: ${GITHUB_TOKEN}"
                  git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.ibm.com/".insteadOf "https://github.ibm.com/"
                '''
              }
              script {
                echo "Switch to branch ${env.BRANCH_NAME}"
                // checkout scm
                sh """
                  go get github.ibm.com/AdvancedAnalyticsCanada/service-config-data
                  ls -ltr  /go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data
                  cd /go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data
                  git checkout ${env.BRANCH_NAME}
                """
              }
            }

            stage('CI Build and Push Snapshot') {
              if ( env.BRANCH_NAME.matches("PR-(.*)") ) {
                script {
                  echo "Building branch ${env.BRANCH_NAME}"
                  sh """
                    cd /go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data
                    make build
                    ls -ltr
                    docker
                    make push
                  """
                }
              }
            }

            stage('Build Release'){
              echo "Building Release"
            }

            stage('Promote to Environments') {
              echo "Promoting to Environment: ${Environments}"
            }

        } catch(error) {
            //echo 'ERROR: ${error}'
        } finally {
            echo "Current BUILD status: ${currentBuild.currentResult} (${currentBuild.displayName})"
        }
        } // container

  } // node
} // pipeline



/*
pipeline {

  agent {
    label "jenkins-go"
  }
  environment {
    ORG = 'AdvancedAnalyticsCanada'
    APP_NAME = 'service-config-data'
    CHARTMUSEUM_CREDS = credentials('jenkins-x-chartmuseum')
  }

  stages {

    stage('Install packages') {
      steps {
      container('go') {
      sh '''
        curl -Lo /tmp/docker.tgz https://get.docker.com/builds/Linux/x86_64/docker-1.12.6.tgz
        mkdir /tmp/docker
        tar -xf /tmp/docker.tgz -C /tmp/docker
        ls -l /tmp/docker/docker/
        mv /tmp/docker/docker/docker* /usr/local/bin/
      '''
      }
      }
    }

    stage('CI Build and push snapshot') {
      when {
        branch 'PR-*'
      }
      environment {
        PREVIEW_VERSION = "0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER"
        PREVIEW_NAMESPACE = "$APP_NAME-$BRANCH_NAME".toLowerCase()
        HELM_RELEASE = "$PREVIEW_NAMESPACE".toLowerCase()
      }
      steps {
        container('go') {

          sh "ls -l /some-secret"
          dir('/home/jenkins/go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data/') {
            checkout scm
            sh "make build"
            sh "make push"
          }
          dir('/home/jenkins/go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data/charts') {
            sh "docker"
            //sh "export VERSION=$PREVIEW_VERSION && skaffold build -f skaffold.yaml"
            //sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:$PREVIEW_VERSION"
          }
        }
      }

    }

    stage('Build Release') {
      when {
        branch 'master'
      }
      steps {
        container('go') {
          dir('/home/jenkins/go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data') {
            checkout scm

            // ensure we're not on a detached head
            sh "git checkout master"
            sh "git config --global credential.helper store"
            sh "jx step git credentials"

            // so we can retrieve the version in later steps
            sh "echo \$(jx-release-version) > VERSION"
            sh "jx step tag --version \$(cat VERSION)"
            sh "make build"
            sh "export VERSION=`cat VERSION` && skaffold build -f skaffold.yaml"
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:\$(cat VERSION)"
          }
        }
      }
    }
    stage('Promote to Environments') {
      when {
        branch 'master'
      }
      steps {
        container('go') {
          dir('/home/jenkins/go/src/github.ibm.com/AdvancedAnalyticsCanada/service-config-data/charts/service-config-data') {
            sh "jx step changelog --version v\$(cat ../../VERSION)"

            // release the helm chart
            sh "jx step helm release"

            // promote through all 'Auto' promotion Environments
            sh "jx promote -b --all-auto --timeout 1h --version \$(cat ../../VERSION)"
          }
        }
      }
    }
  }
}



pipeline {
  agent {
    label "jenkins-go"
  }
  environment {
    ORG = 'AdvancedAnalyticsCanada'
    APP_NAME = 'suncor'
    CHARTMUSEUM_CREDS = credentials('jenkins-x-chartmuseum')
  }

  stages {
    stage('Init stage') {
      steps {
      container('go') {
        echo "Init stage"
      }
      }
    }
  }
}

*/
