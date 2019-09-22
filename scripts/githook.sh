#!/bin/bash
SERVICE_NAME=$1
SERVICE_PORT=$2
ENDPOINT=$3
GIT_USER=$4
GIT_REPO=$5
# get git tocken
TOKEN=$(jq  '.git.access_token' creds.json | tr -d '"')

source ./scripts/common.sh

# ARGS:
# 1 - service name
function wait4external_ip(){
  # waits for LB IP to be provisioned
  external_ip=""
  sp="/-\|"
  echo "Waiting for end point..."
  echo -n ' '
  while [ -z $external_ip ]; do
    echo -ne "\b${sp:i++%${#sp}:1}"
    external_ip=$(kubectl get svc $1 --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
    [ -z "$external_ip" ]  && sleep 1
  done
  #echo -ne "${external_ip}\r"
  echo -ne "\n"
  ret_val=$external_ip
}

printf $BLUE && printf $BRIGHT &&  wait4external_ip $SERVICE_NAME && EXTERNAL_IP=$ret_val && printf $NORMAL
printf $GREEN && printf $BRIGHT && printf $BLINK && echo "Service is running on IP: $EXTERNAL_IP:$SERVICE_PORT "  && printf $NORMAL

# ARGS:
# 1 - Repo
# 2 - user name
# 3 - token
# 4 - IP for hook
function check_githook() {
  # list all githooks
  ret_val=0
  HOOKSLIST=$(curl -s -H "Authorization: token ${3}"  -H "Content-Type: application/json" -X GET https://$GITAPIURL/repos/${2}/${1}/hooks )
  if [ 0 -ne $? ]; then
    printf $BLUE && printf $RED && printf "Something went wrong.. Curl output: " $HOOKSLIST && printf $NORMAL
    ret_val=1
    exit 1
  fi

  printf $BLUE && printf $BRIGHT
  echo "Looking for Githook with IP " ${4} ". Listing existing Githooks:"
  for k in $(jq '.[] | .config.url' <<< "$HOOKSLIST"); do
      echo "k: " $k
      if echo "$k" | grep -q "${4}"; then
        echo "Githook with IP " ${4} "already exists, no need to create";
        ret_val=1
      fi
  done
  printf $NORMAL
}

check_githook $GIT_REPO $GIT_USER $TOKEN $EXTERNAL_IP && githook_flag=$ret_val
if [[ $githook_flag -eq 1 ]] ; then
  exit 0
fi

# ARGS:
# 1 - Repo
# 2 - user name
# 3 - token
# 4 - IP for hook
# 5 - Githook endpoint
# 6 - service port (8000)
function create_githook() {
  # generates json body for githook
  for t in $(find ./ -type f -name "hook.template"); do \
        cat $t | \
          sed -E "s/{{ .WebhookGitSecret }}/$(jq  '.git.git_hook_secret' creds.json | tr -d '"')/g" | \
          sed -E "s/{{ .WebhookGit }}/${5}/g" | \
          sed -E "s/{{ .ServiceIp }}/${4}/g" | \
          sed -E "s/{{ .ServicePort }}/${6}/g"; \
  done  > ./hook.json
  # create githook
  HOOK=$(curl -H "Authorization: token ${3}"  -H "Content-Type: application/json" -vX POST -d @hook.json https://$GITAPIURL/repos/${2}/${1}/hooks)
  if jq -e .type >/dev/null 2>&1 <<<"$HOOK"; then
    printf $GREEN && printf $BRIGHT && echo "Githook created successfully: " $HOOK && printf $NORMAL
    ret_val=0
  else
    printf $RED && printf $BRIGHT && echo "Something went wrong. Message from curl call: " $HOOK && printf $NORMAL
    ret_val=1
  fi
}

create_githook $GIT_REPO $GIT_USER $TOKEN $EXTERNAL_IP $ENDPOINT $SERVICE_PORT  && githook_status=$ret_val
echo $githook_status
