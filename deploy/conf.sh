#!/bin/bash

cat > $CONF <<- EOF
{
    "env": "$ENV",
    "net": {
        "protocol": "http",
        "dns_name": "localhost",
        "port": 8080
    },
    "atlas": {
        "net": {
            "protocol": "https",
            "dns_name": "atlas.ripe.net",
            "port": 443
        },
        "auth":{
            "key_file": "$WORKDIR/atlas.api.key",
            "validate_key": true
        }
    },
    "influxdb": {
        "organization": "dante",
        "net": {
            "protocol": "http",
            "dns_name": "localhost",
            "port": 8086
        },
        "auth": {
            "token_file": "$WORKDIR/influxdb.token"
        }
    },
    "log": {
        "dir": "$LOGDIR"
    }
}
EOF
