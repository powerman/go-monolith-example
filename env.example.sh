#!/bin/sh
# Example default ENV vars for local development.
# Do not modify this file directly, instead, duplicate to (gitignored) `env.sh`.
# Should be loaded into shell used to run `docker-compose up`.

###
### Internal services
###

export MONO_MYSQL_PORT=3306
export MONO_NATS_PORT=4222
export MONO_NATS_URLS="nats://localhost:$MONO_NATS_PORT"
export MONO_STAN_CLUSTER_ID=local

###
### Embedded microservices
###

export MONO_EXAMPLE_DB_USER=root
export MONO_EXAMPLE_DB_PASS=

# DO NOT MODIFY BELOW THIS LINE!
env1="$(sed -e '/^$/d' -e '/^#/d' -e 's/\([A-Z]\)=.*/\1/' env.example.sh)"
env2="$(sed -e '/^$/d' -e '/^#/d' -e 's/\([A-Z]\)=.*/\1/' env.sh)"
if test "$env1" != "$env2"; then
	echo
	echo "[31mFile env.sh differ from env.example.sh, please update and reload env.sh.[0m"
	echo
	return 1
fi
