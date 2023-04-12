#!/bin/sh
exec /usr/bin/rbn-to-kinesis ${STREAM} ${DB_HOST_PORT} ${DB_USER}
