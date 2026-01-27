build:
	go build .

docker-build:
	 docker build --build-arg SERVICE_NAME=bootstrap:latest -t bootstrap:latest .
