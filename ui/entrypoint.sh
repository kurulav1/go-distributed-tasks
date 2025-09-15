#!/bin/sh
set -e
envsubst '${API_BASE}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
nginx -g 'daemon off;'
