build:
	go build .

docker-build:
	 docker build --build-arg SERVICE_NAME=bootstrap VERSION=latest -t bootstrap:latest .
