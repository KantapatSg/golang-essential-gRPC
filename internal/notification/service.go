package notification

import (
	"context"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
	"strings"
	"time"
)

var allowed = map[string]bool{"ORDER_CREATED": true, "ORDER_UPDATED": true, "ORDER_DELETED": true}

type Service struct{ store domain.NotificationStore }

func NewService(store domain.NotificationStore) *Service { return &Service{store: store} }
func (s *Service) Send(ctx context.Context, orderID uuid.UUID, event, message string) (domain.Notification, error) {
	if orderID == uuid.Nil || !allowed[event] || strings.TrimSpace(message) == "" {
		return domain.Notification{}, ErrInvalid
	}
	return s.store.Create(ctx, domain.Notification{ID: uuid.New(), OrderID: orderID, EventType: event, Message: strings.TrimSpace(message), CreatedAt: time.Now().UTC()})
}
func (s *Service) List(ctx context.Context, orderID uuid.UUID, page, size int) ([]domain.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return s.store.List(ctx, orderID, page, size)
}

var ErrInvalid = &invalidError{}

type invalidError struct{}

func (*invalidError) Error() string { return "invalid notification input" }
