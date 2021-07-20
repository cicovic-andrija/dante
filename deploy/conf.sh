#!/bin/bash

cat > $CONF <<- EOF
{
    "env": "$ENV",
    "net": {
        "tcp": {
            "port": $TCP_PORT
        }
    },
    "auth":{
        "key_file": "$WORKDIR/atlas.api.key",
        "validate_key": true
    },
    "log": {
        "dir": "$LOGDIR"
    }
}
EOF
