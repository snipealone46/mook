#!/bin/bash -l

echo "Building the plugin..."
GO111MODULE="on" go build ./cmd/kubectl-mook.go

echo "Installing the plugin..."
rm -rf /usr/local/bin/kubectl-mook
cp ./kubectl-mook /usr/local/bin

echo "Installation complete!!"
