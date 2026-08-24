package service

import (
	"context"
	"testing"

	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
)

type fakeStore struct {
	created []domain.Product
	deleted []uuid.UUID
}

func (f *fakeStore) Create(_ context.Context, p domain.Product) (domain.Product, error) {
	f.created = append(f.created, p)
	return p, nil
}
func (f *fakeStore) Get(context.Context, uuid.UUID) (domain.Product, error) {
	return domain.Product{}, ErrNotFound
}
func (f *fakeStore) List(context.Context, int, int) ([]domain.Product, int64, error) {
	return nil, 0, nil
}
func (f *fakeStore) Update(_ context.Context, p domain.Product) (domain.Product, error) {
	return p, nil
}
func (f *fakeStore) Delete(_ context.Context, id uuid.UUID) error {
	f.deleted = append(f.deleted, id)
	return nil
}

func TestProductServiceCreateTrimsAndGeneratesID(t *testing.T) {
	store := &fakeStore{}
	p, err := NewProductService(store).Create(context.Background(), "  Coffee  ", "  beans ", 12.5, 4)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if p.ID == uuid.Nil || p.Name != "Coffee" || p.Description != "beans" {
		t.Fatalf("unexpected product: %+v", p)
	}
	if len(store.created) != 1 {
		t.Fatalf("created count = %d", len(store.created))
	}
}

func TestProductServiceRejectsInvalidInput(t *testing.T) {
	_, err := NewProductService(&fakeStore{}).Create(context.Background(), "", "", -1, -1)
	if err != ErrInvalidInput {
		t.Fatalf("error = %v, want ErrInvalidInput", err)
	}
}
