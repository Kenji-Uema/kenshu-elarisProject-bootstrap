# Bootstrap

Seeds a fresh environment with shared data and broker topology.

## Main Docs

SSee the main project documentation: <https://kenji-uema.github.io/kenshu-elarisProject-docs/>

## What It Does

- bootstraps MongoDB with cottages, guests, bookings, invoices, receipts, and stock
- declares the RabbitMQ exchanges the application layer expects
- loads cottage photo assets from the mounted photo volume during environment setup

## Runtime Dependencies

- MongoDB
- RabbitMQ
- mounted cottage photo volume
- OpenTelemetry collector

## Local Commands

```sh
go run ./internal
go build ./internal
make docker-build
```

## Minimum Env To Start

Optional vars with defaults, such as `SERVICE_NAME` and `VERSION`, are omitted here.

```sh
PHOTOS_VOLUME_PATH=<mounted photo dir>

MONGO_INITDB_ROOT_USERNAME=<mongo user>
MONGO_INITDB_ROOT_PASSWORD=<mongo password>
MONGO_HOST=<mongo host>
MONGO_DATABASE=cottages
MONGO_COLLECTION_COTTAGE=Cottage
MONGO_COLLECTION_GUEST=Guest
MONGO_COLLECTION_BOOKING=Booking
MONGO_COLLECTION_INVOICE=Invoice
MONGO_COLLECTION_RECEIPT=Receipt
MONGO_COLLECTION_STOCK=Stock

RABBITMQ_USERNAME=<rabbit user>
RABBITMQ_PASSWORD=<rabbit password>
RABBITMQ_HOST=<rabbit host>
RABBITMQ_PORT=5672

CLEANING_EXCHANGE_NAME=ex.cleaning.request
CLEANING_EXCHANGE_KIND=direct
TIME_EVENT_EXCHANGE_NAME=ex.time.event
TIME_EVENT_EXCHANGE_KIND=fanout
INVOICE_EXCHANGE_NAME=ex.invoice.generate
INVOICE_EXCHANGE_KIND=direct
PAYMENT_EXCHANGE_NAME=ex.payment
PAYMENT_EXCHANGE_KIND=topic
COMMUNICATION_EXCHANGE_NAME=ex.communication
COMMUNICATION_EXCHANGE_KIND=direct

OTEL_EXPORTER_OTLP_ENDPOINT=<otel host>
OTEL_EXPORTER_OTLP_GRPC_PORT=4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

## Configuration

Configuration is environment-driven. Start with:

- `internal/config/config.go`
- `internal/config/rabbitmq_config.go`

The main settings groups are service metadata, MongoDB, RabbitMQ, telemetry, and photo/bootstrap resource paths.

## Key Files

- `internal/main.go`
- `internal/infra/mdb/`
- `internal/infra/mq/`
- `resources/`
