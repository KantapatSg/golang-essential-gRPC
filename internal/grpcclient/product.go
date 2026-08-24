package grpcclient

import (
	"context"
	"time"

	productv1 "github.com/KantapatSg/golang-essential-gRPC/gen/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client productv1.ProductServiceClient
}

func NewProductClient(target string) (*ProductClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProductClient{conn: conn, client: productv1.NewProductServiceClient(conn)}, nil
}

func (c *ProductClient) Close() error { return c.conn.Close() }

func (c *ProductClient) Create(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.Product, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return c.client.CreateProduct(ctx, req)
}
func (c *ProductClient) Get(ctx context.Context, req *productv1.GetProductRequest) (*productv1.Product, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return c.client.GetProduct(ctx, req)
}
func (c *ProductClient) List(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return c.client.ListProducts(ctx, req)
}
func (c *ProductClient) Update(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.Product, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	return c.client.UpdateProduct(ctx, req)
}
func (c *ProductClient) Delete(ctx context.Context, req *productv1.DeleteProductRequest) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	_, err := c.client.DeleteProduct(ctx, req)
	return err
}

func (c *ProductClient) Ready(ctx context.Context) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	_, err := c.client.ListProducts(ctx, &productv1.ListProductsRequest{Page: 1, PageSize: 1})
	return err
}

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}
