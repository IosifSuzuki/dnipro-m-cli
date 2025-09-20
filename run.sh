#!/usr/bin/env bash

BINARY_NAME="main"
MAIN_PKG="./internal/cmd"

function main() {
  go build -o $BINARY_NAME $MAIN_PKG
  ./$BINARY_NAME warranty
}


#start point
main