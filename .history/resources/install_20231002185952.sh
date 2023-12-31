#!/bin/sh

go=$(which go)
curl=$(which curl)
dockerCompose=$(which docker-compose)

# copy default env vars
if [ ! -f .env ]; then
    cp .env.dist .env
fi

$dockerCompose build
$dockerCompose up -d

$go get github.com/google/wire/cmd/wire

$go env -w GOPRIVATE=github.com/Montrealist-cPunto

# install linting
$curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.3

# download local
$go mod download
echo ""

$go mod verify

make elasticsearch

# show result
if [ $? -eq 0 ]; then
    echo ""
    echo "Successfully installed!"
    echo ""
else
    echo ""
    echo "Something went wrong. Fix the errors and try again."
    echo ""
fi
