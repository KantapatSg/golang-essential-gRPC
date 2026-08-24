package main

import (
	"context"
	orderv1 "github.com/KantapatSg/golang-essential-gRPC/gen/order"
	"github.com/KantapatSg/golang-essential-gRPC/internal/config"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcclient"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcserver"
	"github.com/KantapatSg/golang-essential-gRPC/internal/repository"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	store, e := repository.NewPostgres(ctx, cfg.OrderDatabaseURL, cfg.DBMaxConns)
	if e != nil {
		log.Fatal(e)
	}
	defer store.Close()
	notifier, e := grpcclient.NewNotificationClient(cfg.NotificationGRPCTarget)
	if e != nil {
		log.Fatal(e)
	}
	defer notifier.Close()
	l, e := net.Listen("tcp", cfg.OrderGRPCAddr)
	if e != nil {
		log.Fatal(e)
	}
	s := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(s, grpcserver.NewOrderServer(service.NewOrderService(store, notifier)))
	go func() { <-ctx.Done(); s.GracefulStop() }()
	log.Printf("order gRPC service listening on %s", cfg.OrderGRPCAddr)
	if e = s.Serve(l); e != nil && ctx.Err() == nil {
		log.Fatal(e)
	}
}
