# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

.PHONY: build
build: binary-build

.PHONY: run
run: build docker-run

#################################
######      Go clean       ######
#################################
.PHONY: clean
clean:

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

#################################
######    Build Binary     ######
#################################
.PHONY: binary-build
binary-build:

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o release/vela-terraform github.com/go-vela/vela-terraform/cmd/vela-terraform

#################################
######    Docker Build     ######
#################################
.PHONY: docker-build
docker-build:

	docker build --no-cache -t terraform-plugin:local .

#################################
######     Docker Run      ######
#################################


current_dir = $(shell pwd)
env = -e PARAMETER_ACTIONS=plan -e GITHUB_TOKEN=some_token
dir = -v ${current_dir}/example:/home/ -w /home/
.PHONY: docker-run
docker-run:

	docker run --rm ${env} ${dir} terraform-plugin:local
