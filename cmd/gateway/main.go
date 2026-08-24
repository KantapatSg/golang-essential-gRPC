package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KantapatSg/golang-essential-gRPC/internal/config"
	"github.com/KantapatSg/golang-essential-gRPC/internal/grpcclient"
	"github.com/KantapatSg/golang-essential-gRPC/internal/httpapi"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()
	client, err := grpcclient.NewProductClient(cfg.GRPCTarget)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	app := fiber.New(fiber.Config{AppName: "product-gateway"})
	httpapi.New(app, client)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() { <-ctx.Done(); _ = app.Shutdown() }()
	log.Printf("Fiber gateway listening on %s; gRPC target %s", cfg.HTTPAddr, cfg.GRPCTarget)
	if err := app.Listen(cfg.HTTPAddr); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}
