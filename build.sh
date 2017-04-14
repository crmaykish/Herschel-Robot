#!/bin/bash

HOSTNAME="herschel.home.cmaykish.com"
SSH_USERNAME="colin"
GO_WS="/home/colin/go/"

echo "Cross-compiling Herschel for Linux on ARM"
env GOOS=linux GOARCH=arm GOARM=7 go build -o "$GO_WS"bin/arm/herschel herschel.go

echo "Copying binary to Herschel..."
scp -q "$GO_WS"bin/arm/herschel "$SSH_USERNAME"@"$HOSTNAME":/home/"$SSH_USERNAME"

echo "Done!"