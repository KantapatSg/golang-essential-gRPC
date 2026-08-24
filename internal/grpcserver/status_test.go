package grpcserver

import (
	"context"

	orderv1 "github.com/KantapatSg/golang-essential-gRPC/gen/order"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

type statusStoreFake struct{}

func (statusStoreFake) Create(_ context.Context, order domain.Order) (domain.Order, error) {
	return order, nil
}
func (statusStoreFake) Get(context.Context, uuid.UUID) (domain.Order, error) {
	return domain.Order{}, service.ErrNotFound
}
func (statusStoreFake) List(context.Context, int, int) ([]domain.Order, int64, error) {
	return nil, 0, nil
}
func (statusStoreFake) Update(_ context.Context, order domain.Order) (domain.Order, error) {
	return order, nil
}
func (statusStoreFake) Delete(context.Context, uuid.UUID) error { return nil }

func TestToStatusMapping(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code codes.Code
	}{{service.ErrInvalidInput, codes.InvalidArgument}, {service.ErrNotFound, codes.NotFound}} {
		if got := status.Code(toStatus(tc.err)); got != tc.code {
			t.Fatalf("got %s want %s", got, tc.code)
		}
	}
}

func TestCreateOrderRejectsUnknownStatus(t *testing.T) {
	server := NewOrderServer(service.NewOrderService(statusStoreFake{}, nil))
	_, err := server.CreateOrder(context.Background(), &orderv1.CreateOrderRequest{
		CustomerName:  "Alice",
		CustomerEmail: "alice@example.com",
		Items:         []*orderv1.OrderItem{{Name: "Book", Quantity: 1, UnitPrice: 10}},
		Status:        orderv1.OrderStatus(-1),
	})
	if got := status.Code(err); got != codes.InvalidArgument {
		t.Fatalf("got %s want %s", got, codes.InvalidArgument)
	}
}
