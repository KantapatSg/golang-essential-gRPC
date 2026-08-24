package service

import (
	"context"
	"errors"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
	"testing"
)

type orderStoreFake struct {
	created domain.Order
	updated domain.Order
	deleted bool
}

func (f *orderStoreFake) Create(_ context.Context, o domain.Order) (domain.Order, error) {
	f.created = o
	return o, nil
}
func (f *orderStoreFake) Get(context.Context, uuid.UUID) (domain.Order, error) {
	return domain.Order{}, nil
}
func (f *orderStoreFake) List(context.Context, int, int) ([]domain.Order, int64, error) {
	return nil, 0, nil
}
func (f *orderStoreFake) Update(_ context.Context, o domain.Order) (domain.Order, error) {
	f.updated = o
	return o, nil
}
func (f *orderStoreFake) Delete(context.Context, uuid.UUID) error { f.deleted = true; return nil }

type notifierFake struct {
	events []string
	err    error
}

func (f *notifierFake) Send(_ context.Context, _ uuid.UUID, event, _ string) error {
	f.events = append(f.events, event)
	return f.err
}
func TestCreateCalculatesTotalAndNotifies(t *testing.T) {
	store := &orderStoreFake{}
	notifier := &notifierFake{}
	s := NewOrderService(store, notifier)
	o, e := s.Create(context.Background(), "Alice", "alice@example.com", []domain.OrderItem{{Name: "Book", Quantity: 2, UnitPrice: 12.5}}, "")
	if e != nil {
		t.Fatal(e)
	}
	if o.TotalAmount != 25 || store.created.TotalAmount != 25 {
		t.Fatalf("total=%v", o.TotalAmount)
	}
	if o.Status != domain.StatusPending || len(notifier.events) != 1 || notifier.events[0] != "ORDER_CREATED" {
		t.Fatalf("status/events: %s %#v", o.Status, notifier.events)
	}
}
func TestInvalidOrderDoesNotStore(t *testing.T) {
	store := &orderStoreFake{}
	s := NewOrderService(store, nil)
	_, e := s.Create(context.Background(), "", "bad", nil, "")
	if !errors.Is(e, ErrInvalidInput) {
		t.Fatalf("err=%v", e)
	}
	if store.created.ID != uuid.Nil {
		t.Fatal("invalid order stored")
	}
}
func TestNotificationFailureDoesNotFailMutation(t *testing.T) {
	store := &orderStoreFake{}
	s := NewOrderService(store, &notifierFake{err: errors.New("down")})
	if _, e := s.Create(context.Background(), "Alice", "alice@example.com", []domain.OrderItem{{Name: "Book", Quantity: 1, UnitPrice: 1}}, ""); e != nil {
		t.Fatalf("notification failure rolled back mutation: %v", e)
	}
}
