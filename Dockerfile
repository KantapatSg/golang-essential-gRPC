FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/product-service ./cmd/product-service && \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/gateway ./cmd/gateway

FROM alpine:3.20 AS product-service
RUN adduser -D -u 10001 app
COPY --from=build /out/product-service /product-service
USER app
EXPOSE 9090
ENTRYPOINT ["/product-service"]

FROM alpine:3.20 AS gateway
RUN adduser -D -u 10001 app
COPY --from=build /out/gateway /gateway
USER app
EXPOSE 8080
ENTRYPOINT ["/gateway"]
