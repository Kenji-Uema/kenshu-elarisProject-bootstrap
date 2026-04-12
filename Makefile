build:
	go build ./internal

docker-build:
	 docker build --build-arg SERVICE_NAME=bootstrap --build-arg VERSION=1.0.0 -t bootstrap:1.0.0 .
