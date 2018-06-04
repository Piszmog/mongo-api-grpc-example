#!/usr/bin/env bash

#protoc -I. --go_out=plugins=micro:movieservice movieservice/movieservice.proto
protoc -I. --micro_out=. --go_out=. movieservice/movieservice.proto
#protoc -I movieservice movieservice/movieservice.proto --go_out=plugins=grpc:movieservice
#go build