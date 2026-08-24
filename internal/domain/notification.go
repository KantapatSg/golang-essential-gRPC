package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	EventType string
	Message   string
	CreatedAt time.Time
}
type NotificationStore interface {
	Create(context.Context, Notification) (Notification, error)
	List(context.Context, uuid.UUID, int, int) ([]Notification, int64, error)
}
