#!/usr/bin/env bash

# This script installs the crontab file located at the root of this repository onto all the desired
# servers via SSH. This allows me to easily modify the cron jobs that are ran on each node in the
# k3s cluster, and makes sure they're all the same.

DIR=$(pwd)
SSH_SERVERS="homelab-0 homelab-1 homelab-2 homelab-3"
SSH_USER="ubuntu"
CRONTAB_DATA=$(cat "${DIR}/crontab")

for SSH_SERVER in $(echo "${SSH_SERVERS}" | sort | uniq); do
  echo "$CRONTAB_DATA" | ssh "${SSH_USER}@${SSH_SERVER}" "crontab -"
done
