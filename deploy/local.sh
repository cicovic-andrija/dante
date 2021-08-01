#!/bin/bash

run_cmd() {
    set -x
    $@ &
    { set +x; } 2>/dev/null
}

WORKDIR="$HOME/dante"
LOGDIR="$WORKDIR/log"
CONF="$WORKDIR/conf.json"
ENV="dev"
TCP_PORT=8080

echo -n "Creating working directory tree $WORKDIR ... "
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
cp deploy/keys/* "$WORKDIR"

echo "Starting dantesrv ..."
run_cmd "$WORKDIR/dantesrv -conf=$CONF"

sleep 1
echo "Done."
