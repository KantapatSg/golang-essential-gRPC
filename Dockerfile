FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/gateway ./cmd/gateway && CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/order-service ./cmd/order-service && CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/notification-service ./cmd/notification-service
FROM alpine:3.20 AS gateway
RUN adduser -D -u 10001 app
COPY --from=build /out/gateway /gateway
USER app
EXPOSE 8080
ENTRYPOINT ["/gateway"]
FROM alpine:3.20 AS order-service
RUN adduser -D -u 10001 app
COPY --from=build /out/order-service /order-service
USER app
EXPOSE 9090
ENTRYPOINT ["/order-service"]
FROM alpine:3.20 AS notification-service
RUN adduser -D -u 10001 app
COPY --from=build /out/notification-service /notification-service
USER app
EXPOSE 9091
ENTRYPOINT ["/notification-service"]
