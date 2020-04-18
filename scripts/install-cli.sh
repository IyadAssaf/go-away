#!/bin/bash

set -e

GO111MODULE=on
INSTALL_PATH=/usr/local/bin/
echo Installing in $INSTALL_PATH
go build ./cmd/go-away
mv ./go-away $INSTALL_PATH
chmod +x $INSTALL_PATH/go-away