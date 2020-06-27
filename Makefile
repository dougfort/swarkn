SHELL := /bin/bash

export PROJECT = swiss-army-knife

# Building containers

all: server client

server:
	docker build \
		-f dockerfile.server \
		-t deciphernow/swiss_army_knife_server:latest \
		--build-arg PACKAGE_NAME=swarkn \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

client:
	docker build \
		-f dockerfile.client \
		-t deciphernow/swiss_army_knife_client:latest \
		--build-arg PACKAGE_NAME=swarkn \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

