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
        "net": {
            "protocol": "http",
            "dns_name": "localhost",
            "port": 8086
        }
    },
    "log": {
        "dir": "$LOGDIR"
    }
}
EOF
