package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	Name      string
	Quantity  int32
	UnitPrice float64
}
type Order struct {
	ID            uuid.UUID
	CustomerName  string
	CustomerEmail string
	Items         []OrderItem
	TotalAmount   float64
	Status        OrderStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderStore interface {
	Create(context.Context, Order) (Order, error)
	Get(context.Context, uuid.UUID) (Order, error)
	List(context.Context, int, int) ([]Order, int64, error)
	Update(context.Context, Order) (Order, error)
	Delete(context.Context, uuid.UUID) error
}

type Notifier interface {
	Send(context.Context, uuid.UUID, string, string) error
}
