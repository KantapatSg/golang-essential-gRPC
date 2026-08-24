package main

import (
	"context"
	"github.com/KantapatSg/golang-essential-gRPC/internal/config"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcclient"
	"github.com/KantapatSg/golang-essential-gRPC/internal/httpapi"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	orders, e := grpcclient.NewOrderClient(cfg.OrderGRPCTarget)
	if e != nil {
		log.Fatal(e)
	}
	defer orders.Close()
	notifications, e := grpcclient.NewNotificationClient(cfg.NotificationGRPCTarget)
	if e != nil {
		log.Fatal(e)
	}
	defer notifications.Close()
	app := fiber.New(fiber.Config{AppName: "order-gateway"})
	httpapi.New(app, orders, notifications)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() { <-ctx.Done(); _ = app.Shutdown() }()
	log.Printf("Fiber gateway listening on %s", cfg.HTTPAddr)
	if e := app.Listen(cfg.HTTPAddr); e != nil && ctx.Err() == nil {
		log.Fatal(e)
	}
}
