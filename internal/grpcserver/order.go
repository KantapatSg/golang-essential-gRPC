package grpcserver

import (
	"context"
	"errors"
	orderv1 "github.com/KantapatSg/golang-essential-gRPC/gen/order"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	svc *service.OrderService
}

func NewOrderServer(s *service.OrderService) *OrderServer { return &OrderServer{svc: s} }
func (s *OrderServer) CreateOrder(ctx context.Context, r *orderv1.CreateOrderRequest) (*orderv1.Order, error) {
	o, e := s.svc.Create(ctx, r.GetCustomerName(), r.GetCustomerEmail(), fromItems(r.GetItems()), fromStatus(r.GetStatus()))
	if e != nil {
		return nil, toStatus(e)
	}
	return toOrder(o), nil
}
func (s *OrderServer) GetOrder(ctx context.Context, r *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	id, e := uuid.Parse(r.GetId())
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	o, e := s.svc.Get(ctx, id)
	if e != nil {
		return nil, toStatus(e)
	}
	return toOrder(o), nil
}
func (s *OrderServer) ListOrders(ctx context.Context, r *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	items, total, e := s.svc.List(ctx, int(r.GetPage()), int(r.GetPageSize()))
	if e != nil {
		return nil, toStatus(e)
	}
	page, size := r.GetPage(), r.GetPageSize()
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	out := &orderv1.ListOrdersResponse{Page: page, PageSize: size, Total: total, Orders: make([]*orderv1.Order, 0, len(items))}
	for _, o := range items {
		out.Orders = append(out.Orders, toOrder(o))
	}
	return out, nil
}
func (s *OrderServer) UpdateOrder(ctx context.Context, r *orderv1.UpdateOrderRequest) (*orderv1.Order, error) {
	id, e := uuid.Parse(r.GetId())
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	o, e := s.svc.Update(ctx, id, r.GetCustomerName(), r.GetCustomerEmail(), fromItems(r.GetItems()), fromStatus(r.GetStatus()))
	if e != nil {
		return nil, toStatus(e)
	}
	return toOrder(o), nil
}
func (s *OrderServer) DeleteOrder(ctx context.Context, r *orderv1.DeleteOrderRequest) (*emptypb.Empty, error) {
	id, e := uuid.Parse(r.GetId())
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, "id must be a valid UUID")
	}
	if e = s.svc.Delete(ctx, id); e != nil {
		return nil, toStatus(e)
	}
	return &emptypb.Empty{}, nil
}
func fromItems(in []*orderv1.OrderItem) []domain.OrderItem {
	out := make([]domain.OrderItem, len(in))
	for i, v := range in {
		out[i] = domain.OrderItem{Name: v.GetName(), Quantity: v.GetQuantity(), UnitPrice: v.GetUnitPrice()}
	}
	return out
}
func fromStatus(v orderv1.OrderStatus) domain.OrderStatus {
	switch v {
	case orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED:
		return ""
	case orderv1.OrderStatus_PENDING:
		return domain.StatusPending
	case orderv1.OrderStatus_CONFIRMED:
		return domain.StatusConfirmed
	case orderv1.OrderStatus_CANCELLED:
		return domain.StatusCancelled
	default:
		return domain.OrderStatus("INVALID")
	}
}
func toOrder(o domain.Order) *orderv1.Order {
	items := make([]*orderv1.OrderItem, len(o.Items))
	for i, v := range o.Items {
		items[i] = &orderv1.OrderItem{Name: v.Name, Quantity: v.Quantity, UnitPrice: v.UnitPrice}
	}
	var st orderv1.OrderStatus
	switch o.Status {
	case domain.StatusConfirmed:
		st = orderv1.OrderStatus_CONFIRMED
	case domain.StatusCancelled:
		st = orderv1.OrderStatus_CANCELLED
	default:
		st = orderv1.OrderStatus_PENDING
	}
	return &orderv1.Order{Id: o.ID.String(), CustomerName: o.CustomerName, CustomerEmail: o.CustomerEmail, Items: items, TotalAmount: o.TotalAmount, Status: st, CreatedAt: timestamppb.New(o.CreatedAt), UpdatedAt: timestamppb.New(o.UpdatedAt)}
}
func toStatus(e error) error {
	switch {
	case errors.Is(e, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, e.Error())
	case errors.Is(e, service.ErrNotFound):
		return status.Error(codes.NotFound, e.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
