FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bootstrap .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bootstrap /app/bootstrap
COPY --from=builder /app/resources /app/resources

ENTRYPOINT ["/app/bootstrap"]
