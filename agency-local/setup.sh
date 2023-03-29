#!/bin/bash

set -e

CURRENT_DIR=$(dirname "$BASH_SOURCE")

curl -s https://raw.githubusercontent.com/findy-network/e2e-test-action/master/env/cert/client/client-pkcs8.key >$CURRENT_DIR/cert/client/client-pkcs8.key
curl -s https://raw.githubusercontent.com/findy-network/e2e-test-action/master/env/cert/client/client.key >$CURRENT_DIR/cert/client/client.key
curl -s https://raw.githubusercontent.com/findy-network/e2e-test-action/master/env/cert/client/client.crt >$CURRENT_DIR/cert/client/client.crt
curl -s https://raw.githubusercontent.com/findy-network/e2e-test-action/master/env/cert/server/server.key >$CURRENT_DIR/cert/server/server.key
curl -s https://raw.githubusercontent.com/findy-network/e2e-test-action/master/env/cert/server/server.crt >$CURRENT_DIR/cert/server/server.crt

docker-compose -f $CURRENT_DIR/docker-compose.yml up --pull="always"
