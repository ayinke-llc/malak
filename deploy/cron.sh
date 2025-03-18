#!/bin/bash

## fetch the token from Infisical
## safe to store here. token Only valid for about 60 seconds

# Check if required environment variables are set
if [ -z "$INFISICAL_CLIENT_ID" ]; then
	echo "Error: INFISICAL_CLIENT_ID is not set"
	exit 1
fi

if [ -z "$INFISICAL_CLIENT_SECRET" ]; then
	echo "Error: INFISICAL_CLIENT_SECRET is not set"
	exit 1
fi

mkdir -p /etc/infisical
export INFISICAL_TOKEN=$(infisical login --method=universal-auth --client-id="${INFISICAL_CLIENT_ID}" --client-secret="${INFISICAL_CLIENT_SECRET}" --silent --plain)

if [ -z "$INFISICAL_TOKEN" ]; then
	echo "Error: Failed to obtain INFISICAL_TOKEN"
	exit 1
fi

sudo echo "INFISICAL_TOKEN=$INFISICAL_TOKEN" >/etc/infisical/infisical.env

sudo chmod 600 /etc/infisical/infisical.env

infisical run --projectId f7d168e1-36bc-4fc4-9a0f-3b26d837187e --env=prod -- malak cron $CRON_TYPE
