#!/bin/sh
# Example default ENV vars for local development.
# Do not modify `env.sh.dist` directly, copy it to (gitignored) `env.sh` and use that instead.
# Should be loaded into shell used to run `docker-compose up`.

# - Set all _PORT vars to port numbers not used by your system.

export EASYRSA_PKI=configs/insecure-dev-pki
export GO_TEST_TIME_FACTOR="1.0" # Increase if tests fail because of slow CPU.

# Lower-case variables are either used only by docker-compose.yml or
# provide reusable values for project's upper-case variables defined below.
export mono_mysql_addr_port="3306"
export mono_nats_addr_port="4222"
export mono_postgres_addr_port="5432"
export mono_postgres_tls_cert="./$EASYRSA_PKI/issued/postgres.crt"
export mono_postgres_tls_key="./$EASYRSA_PKI/private/postgres.key"

# Variables required to run and test project.
# Should be kept in sorted order.
# Avoid referencing one variable from another if their order may change,
# use lower-case variables defined above for such a shared values.
# Naming convention:
#   <PROJECT>_<VAR>         - global vars, not specific for some embedded microservice (e.g. domain)
#   <PROJECT>_X_<SVC>_<VAR> - vars related to external services (e.g. databases)
#   <PROJECT>_<MS>_<VAR>    - vars related to embedded microservice (e.g. addr)
#   <PROJECT>__<MS>_<VAR>   - private vars for embedded microservice
export MONO_ADDR_HOST="localhost"
export MONO_ADDR_HOST_INT="127.0.0.1"
export MONO_AUTH_ADDR_HOST="localhost"     # Must match DNS/IP in $MONO__AUTH_TLS_CERT.
export MONO_AUTH_ADDR_HOST_INT="127.0.0.1" # Must match DNS/IP in $MONO__AUTH_TLS_CERT_INT.
export MONO_AUTH_ADDR_PORT="17003"
export MONO_AUTH_ADDR_PORT_INT="17004"
export MONO_AUTH_GRPCGW_ADDR_PORT="17006"
export MONO_AUTH_METRICS_ADDR_PORT="17005"
export MONO_EXAMPLE_ADDR_PORT="17001"
export MONO_EXAMPLE_METRICS_ADDR_PORT="17002"
export MONO_MONO_ADDR_PORT="17000"
export MONO_TLS_CA_CERT="./$EASYRSA_PKI/ca.crt"
export MONO_X_MYSQL_ADDR_HOST="localhost"
export MONO_X_MYSQL_ADDR_PORT="${mono_mysql_addr_port}"
export MONO_X_NATS_ADDR_URLS="nats://localhost:${mono_nats_addr_port}"
export MONO_X_POSTGRES_ADDR_HOST="localhost"
export MONO_X_POSTGRES_ADDR_PORT="${mono_postgres_addr_port}"
export MONO_X_STAN_CLUSTER_ID="local"
export MONO__AUTH_POSTGRES_AUTH_LOGIN="auth"
export MONO__AUTH_POSTGRES_AUTH_PASS="authpass"
export MONO__AUTH_SECRET="s3cr3t"
export MONO__AUTH_TLS_CERT="./$EASYRSA_PKI/issued/ms-auth.crt"
export MONO__AUTH_TLS_CERT_INT="./$EASYRSA_PKI/issued/ms-auth-int.crt"
export MONO__AUTH_TLS_KEY="./$EASYRSA_PKI/private/ms-auth.key"
export MONO__AUTH_TLS_KEY_INT="./$EASYRSA_PKI/private/ms-auth-int.key"
export MONO__EXAMPLE_MYSQL_AUTH_LOGIN="root"
export MONO__EXAMPLE_MYSQL_AUTH_PASS=""

# DO NOT MODIFY BELOW THIS LINE!
env1="$(sed -e '/^$/d' -e '/^#/d' -e 's/=.*//' env.sh.dist)"
env2="$(sed -e '/^$/d' -e '/^#/d' -e 's/=.*//' env.sh)"
if test "$env1" != "$env2"; then
	echo
	echo "[31mFile env.sh differ from env.sh.dist, please update and reload env.sh.[0m"
	echo
	return 1
fi
