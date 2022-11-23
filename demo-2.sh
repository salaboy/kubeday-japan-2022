#!/usr/bin/env bash

########################
# include the magic
########################
. demo-lib/demo-magic.sh


########################
# Configure the options
########################

#
# speed at which to simulate typing. bigger num = faster
#
# TYPE_SPEED=20

#
# custom prompt
#
# see http://www.tldp.org/HOWTO/Bash-Prompt-HOWTO/bash-prompt-escape-sequences.html for escape sequences
#

# hide the evidence
clear

pe "cd functions"
pe "mkdir avg && cd avg"
pe "func create -r https://github.com/salaboy/func -l go -t redis"
pe "ls -al"
pe "code ."
pe "func config envs add --name=REDIS_HOST --value='kubeday-japan-app-redis-master-x-default-x-team-a-env:6379'"
pe "func config envs add --name=REDIS_PASSWORD --value='{{ secret:kubeday-japan-app-redis-x-default-x-team-a-env:redis-password }}'"
pe "func deploy -v --registry docker.io/salaboy"



