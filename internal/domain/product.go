package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductStore interface {
	Create(context.Context, Product) (Product, error)
	Get(context.Context, uuid.UUID) (Product, error)
	List(context.Context, int, int) ([]Product, int64, error)
	Update(context.Context, Product) (Product, error)
	Delete(context.Context, uuid.UUID) error
}
