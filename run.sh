#!/bin/sh

pushd client
yarn
rm dist/*
yarn build
popd

go build .
PORT=80 GIN_MODE=release ./fela
