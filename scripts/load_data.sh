#! /bin/bash

docker-compose exec covid-server go run ./cmd/loader/main.go $@
