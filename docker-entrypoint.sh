#!/usr/bin/env bash
set -Eeo pipefail
# TODO swap to -Eeuo pipefail above (after handling all potentially-unset variables)
# usage: file_env VAR [DEFAULT]
#    ie: file_env 'XYZ_DB_PASSWORD' 'example'
# (will allow for "$XYZ_DB_PASSWORD_FILE" to fill in the value of
#  "$XYZ_DB_PASSWORD" from a file, especially for Docker's secrets feature)
file_env() {
	local var="$1"
	local fileVar="${var}_FILE"
	local def="${2:-}"
	if [ "${!var:-}" ] && [ "${!fileVar:-}" ]; then
		echo >&2 "error: both $var and $fileVar are set (but are exclusive)"
		exit 1
	fi
	local val="$def"
	if [ "${!var:-}" ]; then
		val="${!var}"
	elif [ "${!fileVar:-}" ]; then
		val="$(< "${!fileVar}")"
	fi
	export "$var"="$val"
	unset "$fileVar"
}

# check to see if this file is being run or sourced from another script
_is_sourced() {
	# https://unix.stackexchange.com/a/215279
	[ "${#FUNCNAME[@]}" -ge 2 ] \
		&& [ "${FUNCNAME[0]}" = '_is_sourced' ] \
		&& [ "${FUNCNAME[1]}" = 'source' ]
}


# Loads various settings that are used elsewhere in the script
# This should be called before any other functions
docker_setup_env() {
    file_env 'PHOENIX_DATABASE__HOST'
	file_env 'PHOENIX_DATABASE__PASSWORD'
	file_env 'PHOENIX_DATABASE__USER'
    file_env 'PHOENIX_DATABASE__NAME'
	file_env 'POSTGRES_INITDB_ARGS'
}

docker_run_pgsd(){
    exec pqsd -connect "postgresql://postgres:pass@localhost/iot?sslmode=disable"
}

_main() {
    docker_setup_env
    docker_run_pgsd
}

if ! _is_sourced; then
	_main 
fi