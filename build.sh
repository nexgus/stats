#!/bin/bash
PROJ=stats
if [ "$#" -eq 0 ]; then
    TARGET=${PROJ}
else
    TARGET=$1
fi

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags '-w -extldflags "-static"' \
    -o bin/$TARGET \
    ${PROJ}/cmd/${TARGET}
