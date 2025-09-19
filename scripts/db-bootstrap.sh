#!/usr/bin/env bash
set -euo pipefail

RPC_HOST="${RPC_HOST:-roach1}"
RPC_PORT="${RPC_PORT:-26257}"
SQL_HOST="${SQL_HOST:-roach1}"
SQL_PORT="${SQL_PORT:-26258}"
DB_NAME="${DB_NAME:-app}"
SCHEMA_FILE="${SCHEMA_FILE:-/schema.sql}"

until /cockroach/cockroach sql --insecure --host="${SQL_HOST}:${SQL_PORT}" -e "select 1" >/dev/null 2>&1; do
  sleep 1
done

/cockroach/cockroach init --insecure --cluster-name=distrib --host="${RPC_HOST}:${RPC_PORT}" || true

/cockroach/cockroach sql --insecure --host="${SQL_HOST}:${SQL_PORT}" -e "create database if not exists ${DB_NAME}"

if [ -f "${SCHEMA_FILE}" ]; then
  /cockroach/cockroach sql --insecure --host="${SQL_HOST}:${SQL_PORT}" -d "${DB_NAME}" -f "${SCHEMA_FILE}"
fi

/cockroach/cockroach sql --insecure --host="${SQL_HOST}:${SQL_PORT}" -d "${DB_NAME}" -e "show tables"
