package service

import (
	"context"
	"strings"
	"time"

	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type ProductService struct{ store domain.ProductStore }

func NewProductService(store domain.ProductStore) *ProductService {
	return &ProductService{store: store}
}

func (s *ProductService) Create(ctx context.Context, name, description string, price float64, stock int32) (domain.Product, error) {
	if err := validate(name, price, stock); err != nil {
		return domain.Product{}, err
	}
	now := time.Now().UTC()
	return s.store.Create(ctx, domain.Product{ID: uuid.New(), Name: strings.TrimSpace(name), Description: strings.TrimSpace(description), Price: price, Stock: stock, CreatedAt: now, UpdatedAt: now})
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return s.store.Get(ctx, id)
}

func (s *ProductService) List(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	if page < 1 {
		page = defaultPage
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return s.store.List(ctx, page, pageSize)
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, name, description string, price float64, stock int32) (domain.Product, error) {
	if err := validate(name, price, stock); err != nil {
		return domain.Product{}, err
	}
	product := domain.Product{ID: id, Name: strings.TrimSpace(name), Description: strings.TrimSpace(description), Price: price, Stock: stock, UpdatedAt: time.Now().UTC()}
	return s.store.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.Delete(ctx, id)
}

func validate(name string, price float64, stock int32) error {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 120 || price < 0 || stock < 0 {
		return ErrInvalidInput
	}
	return nil
}
