FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o upassed-assignment-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /upassed-assignment-service/config
RUN mkdir -p /upassed-assignment-service/migration/scripts
COPY --from=builder /app/upassed-assignment-service /upassed-assignment-service/upassed-assignment-service
COPY --from=builder /app/config/* /upassed-assignment-service/config
COPY --from=builder /app/migration/scripts/* /upassed-assignment-service/migration/scripts
RUN chmod +x /upassed-assignment-service/upassed-assignment-service
ENV APP_CONFIG_PATH="/upassed-assignment-service/config/local.yml"
EXPOSE 44047
CMD ["/upassed-assignment-service/upassed-assignment-service"]
