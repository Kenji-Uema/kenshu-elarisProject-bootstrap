build:
	go build .

docker-build:
	 docker build --build-arg SERVICE_NAME=bootstrap --build-arg VERSION=latest -t bootstrap:latest .
