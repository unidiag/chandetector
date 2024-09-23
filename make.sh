#!/bin/sh

rm ./chandetector
go build -ldflags "-linkmode external -extldflags '-static'" -o chandetector
