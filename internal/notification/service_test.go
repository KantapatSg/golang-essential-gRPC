package notification

import (
	"context"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
	"testing"
)

type storeFake struct{ n domain.Notification }

func (f *storeFake) Create(_ context.Context, n domain.Notification) (domain.Notification, error) {
	f.n = n
	return n, nil
}
func (f *storeFake) List(context.Context, uuid.UUID, int, int) ([]domain.Notification, int64, error) {
	return nil, 0, nil
}
func TestSendPersistsAllowedEvent(t *testing.T) {
	f := &storeFake{}
	s := NewService(f)
	id := uuid.New()
	if _, e := s.Send(context.Background(), id, "ORDER_CREATED", "created"); e != nil {
		t.Fatal(e)
	}
	if f.n.OrderID != id || f.n.EventType != "ORDER_CREATED" {
		t.Fatalf("%+v", f.n)
	}
}
