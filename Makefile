.PHONY: proto tidy fmt test vet compose-up compose-down
PROTOC_INCLUDE?=
proto:
	protoc -I proto $(if $(PROTOC_INCLUDE),-I "$(PROTOC_INCLUDE)") --go_out=. --go_opt=module=github.com/KantapatSg/golang-essential-gRPC --go-grpc_out=. --go-grpc_opt=module=github.com/KantapatSg/golang-essential-gRPC proto/order.proto proto/notification.proto
tidy:
	go mod tidy
fmt:
	go fmt ./...
test:
	go test ./...
vet:
	go vet ./...
compose-up:
	docker compose up --build
compose-down:
	docker compose down
