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

