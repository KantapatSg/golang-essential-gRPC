package main

import (
	"context"
	notificationv1 "github.com/KantapatSg/golang-essential-gRPC/gen/notification"
	"github.com/KantapatSg/golang-essential-gRPC/internal/config"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcserver"
	"github.com/KantapatSg/golang-essential-gRPC/internal/notification"
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
	store, e := notification.NewPostgres(ctx, cfg.NotificationDatabaseURL, cfg.DBMaxConns)
	if e != nil {
		log.Fatal(e)
	}
	defer store.Close()
	l, e := net.Listen("tcp", cfg.NotificationGRPCAddr)
	if e != nil {
		log.Fatal(e)
	}
	s := grpc.NewServer()
	notificationv1.RegisterNotificationServiceServer(s, grpcserver.NewNotificationServer(notification.NewService(store)))
	go func() { <-ctx.Done(); s.GracefulStop() }()
	log.Printf("notification gRPC service listening on %s", cfg.NotificationGRPCAddr)
	if e = s.Serve(l); e != nil && ctx.Err() == nil {
		log.Fatal(e)
	}
}
