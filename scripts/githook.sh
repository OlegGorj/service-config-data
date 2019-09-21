#!/bin/bash

SERVICE_NAME=$1
SERVICE_PORT=$2
ENDPOINT=$3
GIT_USER=$4
GIT_REPO=$5

# waits for LB IP to be provisioned
external_ip=""
sp="/-\|"
echo "Waiting for end point..."
echo -n ' '
while [ -z $external_ip ]; do
  echo -ne "\b${sp:i++%${#sp}:1}"
  external_ip=$(kubectl get svc $SERVICE_NAME --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
  [ -z "$external_ip" ] # && sleep 10
done
echo -ne "${external_ip}\r"
echo -ne "\n"

# generates json body for githook
for t in $(find ./ -type f -name "hook.template"); do \
      cat $t | \
        sed -E "s/{{ .WebhookGitSecret }}/$(jq  '.git.git_hook_secret' creds.json | tr -d '"')/g" | \
        sed -E "s/{{ .WebhookGit }}/${ENDPOINT}/g" | \
        sed -E "s/{{ .ServiceIp }}/${external_ip}/g" | \
        sed -E "s/{{ .ServicePort }}/${SERVICE_PORT}/g"; \
done  > ./hook.json

# get git tocken
TOKEN=$(jq  '.git.access_token' creds.json | tr -d '"')

# create githook
HOOK=$(curl -H "Authorization: token $TOKEN"  -H "Content-Type: application/json" -vX POST -d @hook.json https://api.github.com/repos/${GIT_USER}/${GIT_REPO}/hooks)
echo "Curl Response: " $HOOK

if jq -e .type >/dev/null 2>&1 <<<"$HOOK"; then
  echo "Hook created successfully"
else
  echo "Something went wrong"
fi
