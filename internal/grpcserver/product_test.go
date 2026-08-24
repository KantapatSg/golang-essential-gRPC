package grpcserver

import (
	"context"
	"testing"

	productv1 "github.com/KantapatSg/golang-essential-gRPC/gen/product"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeStore struct{}

func (fakeStore) Create(_ context.Context, p domain.Product) (domain.Product, error) { return p, nil }
func (fakeStore) Get(context.Context, uuid.UUID) (domain.Product, error) {
	return domain.Product{}, service.ErrNotFound
}
func (fakeStore) List(context.Context, int, int) ([]domain.Product, int64, error)    { return nil, 0, nil }
func (fakeStore) Update(_ context.Context, p domain.Product) (domain.Product, error) { return p, nil }
func (fakeStore) Delete(context.Context, uuid.UUID) error                            { return service.ErrNotFound }

func TestGetProductMapsInvalidUUID(t *testing.T) {
	_, err := NewProductServer(service.NewProductService(fakeStore{})).GetProduct(context.Background(), &productv1.GetProductRequest{Id: "bad-id"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status = %v, want InvalidArgument", status.Code(err))
	}
}

func TestGetProductMapsNotFound(t *testing.T) {
	_, err := NewProductServer(service.NewProductService(fakeStore{})).GetProduct(context.Background(), &productv1.GetProductRequest{Id: uuid.NewString()})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("status = %v, want NotFound", status.Code(err))
	}
}
