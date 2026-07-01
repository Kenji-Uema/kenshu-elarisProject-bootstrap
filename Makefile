IMAGE_TAG ?= 1.0.0

build:
	go build ./internal

docker-build:
	 docker buildx build --build-arg SERVICE_NAME=bootstrap --build-arg VERSION=$(IMAGE_TAG) -t bootstrap:$(IMAGE_TAG) --load .
