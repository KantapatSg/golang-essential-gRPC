package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	productv1 "github.com/KantapatSg/golang-essential-gRPC/gen/product"
	"github.com/KantapatSg/golang-essential-gRPC/internal/config"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcserver"
	"github.com/KantapatSg/golang-essential-gRPC/internal/repository"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	store, err := repository.NewPostgres(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	productv1.RegisterProductServiceServer(server, grpcserver.NewProductServer(service.NewProductService(store)))
	go func() { <-ctx.Done(); server.GracefulStop() }()
	log.Printf("product gRPC service listening on %s", cfg.GRPCAddr)
	if err := server.Serve(listener); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}
