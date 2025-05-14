#!/bin/bash

get_args() {
	_arg=$1
	shift
	echo "$@" | grep -Eo "\-\-${_arg}=[^ ]+" | cut -d= -f2 | tail -n 1
	unset _arg
}

subscriptions=$(get_args subscriptions "$@")

configPath="./settings/subscriptions.json"

[ "$GRAVITY_EXPORTER_JETSTREAM_RULES_SUBSCRIPTION" != "" ] && {
	configPath=$GRAVITY_EXPORTER_JETSTREAM_RULES_SUBSCRIPTION
} 

[ "$subscriptions" != "" ] && {
	echo $subscriptions > $configPath
}

[ "$subscriptions" != "" ] || {
	echo $GRAVITY_EXPORTER_JETSTREAM_SUBSCRIPTION_SETTINGS > $configPath
}

exec /gravity-exporter-jetstream
