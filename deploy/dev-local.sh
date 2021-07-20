#!/bin/bash

WORKDIR="$HOME/dante"
LOGDIR="$WORKDIR/log"
CONF="$WORKDIR/conf.json"
ENV="dev"
TCP_PORT=8080

echo -n "Creating working directory $WORKDIR ... "
mkdir -p "$WORKDIR"
mkdir -p "$LOGDIR"
echo "Done."

echo -n "Writing config to $CONF ... "
source deploy/conf.sh
echo "Done."

echo -n "Building project ... "
GOBIN="$WORKDIR" go install cmd/dantesrv.go
if [ $? -eq 0 ]; then
    echo "Done."
else
    exit 1
fi

# Copy secrets
source deploy/dev-copy-secrets.sh

echo "Starting dantesrv ..."
set -x
$WORKDIR/dantesrv -conf="$CONF"
{ set +x; } 2>/dev/null

sleep 1
echo "Done."
