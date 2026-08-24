package grpcserver

import (
	"context"
	"errors"
	notificationv1 "github.com/KantapatSg/golang-essential-gRPC/gen/notification"
	"github.com/KantapatSg/golang-essential-gRPC/internal/notification"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	svc *notification.Service
}

func NewNotificationServer(s *notification.Service) *NotificationServer {
	return &NotificationServer{svc: s}
}
func (s *NotificationServer) SendNotification(ctx context.Context, r *notificationv1.SendNotificationRequest) (*notificationv1.Notification, error) {
	id, e := uuid.Parse(r.GetOrderId())
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, "order_id must be a valid UUID")
	}
	n, e := s.svc.Send(ctx, id, r.GetEventType(), r.GetMessage())
	if e != nil {
		if errors.Is(e, notification.ErrInvalid) {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &notificationv1.Notification{Id: n.ID.String(), OrderId: n.OrderID.String(), EventType: n.EventType, Message: n.Message, CreatedAt: timestamppb.New(n.CreatedAt)}, nil
}
func (s *NotificationServer) ListNotifications(ctx context.Context, r *notificationv1.ListNotificationsRequest) (*notificationv1.ListNotificationsResponse, error) {
	id := uuid.Nil
	var e error
	if r.GetOrderId() != "" {
		id, e = uuid.Parse(r.GetOrderId())
		if e != nil {
			return nil, status.Error(codes.InvalidArgument, "order_id must be a valid UUID")
		}
	}
	items, total, e := s.svc.List(ctx, id, int(r.GetPage()), int(r.GetPageSize()))
	if e != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	p, size := r.GetPage(), r.GetPageSize()
	if p < 1 {
		p = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	out := &notificationv1.ListNotificationsResponse{Page: p, PageSize: size, Total: total, Notifications: make([]*notificationv1.Notification, 0, len(items))}
	for _, n := range items {
		out.Notifications = append(out.Notifications, &notificationv1.Notification{Id: n.ID.String(), OrderId: n.OrderID.String(), EventType: n.EventType, Message: n.Message, CreatedAt: timestamppb.New(n.CreatedAt)})
	}
	return out, nil
}
