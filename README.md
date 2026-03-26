# bootstrap

Seeds the platform with initial data and exchange definitions.

## Responsibilities

- bootstrap MongoDB collections with cottages, guests, bookings, invoices, receipts, and stock
- bootstrap RabbitMQ exchanges used by the application layer
- prepare shared photo assets through the mounted volume

## Depends on

- MongoDB
- RabbitMQ
- OpenTelemetry collector
- mounted photo volume

## Run

```sh
go run ./internal
```

## Build

```sh
make build
make docker-build
```

## Configuration

Configuration is environment-driven. See:

- `internal/config/config.go`
- `internal/config/rabbitmq_config.go`

## Entry point

- `internal/main.go`
