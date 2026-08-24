APP?=product-api
PROTOC_INCLUDE?=

.PHONY: proto tidy fmt test vet run-api run-service compose-up compose-down

proto:
	protoc -I proto $(if $(PROTOC_INCLUDE),-I "$(PROTOC_INCLUDE)") --go_out=. --go_opt=module=github.com/KantapatSg/golang-essential-gRPC --go-grpc_out=. --go-grpc_opt=module=github.com/KantapatSg/golang-essential-gRPC proto/product.proto

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

run-api:
	go run ./cmd/gateway

run-service:
	go run ./cmd/product-service

compose-up:
	docker compose up --build

compose-down:
	docker compose down
