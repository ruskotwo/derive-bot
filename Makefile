TAG=ruskotwo/derive-bot:latest
DOCKER_BUILD_OPTIONS ?=--platform linux/amd64

wire:
	cd cmd/factory && wire ; cd ../..

golang_build:
	docker build $(DOCKER_BUILD_OPTIONS) \
		-t $(TAG) -f ./docker/golang.Dockerfile .
