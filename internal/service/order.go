package service

import (
	"context"
	"fmt"
	"net/mail"
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

type OrderService struct {
	store    domain.OrderStore
	notifier domain.Notifier
}

func NewOrderService(store domain.OrderStore, notifier domain.Notifier) *OrderService {
	return &OrderService{store: store, notifier: notifier}
}

func (s *OrderService) Create(ctx context.Context, name, email string, items []domain.OrderItem, status domain.OrderStatus) (domain.Order, error) {
	cleanName, cleanEmail, cleanItems, total, err := validate(name, email, items, status)
	if err != nil {
		return domain.Order{}, err
	}
	if status == "" {
		status = domain.StatusPending
	}
	now := time.Now().UTC()
	o, err := s.store.Create(ctx, domain.Order{ID: uuid.New(), CustomerName: cleanName, CustomerEmail: cleanEmail, Items: cleanItems, TotalAmount: total, Status: status, CreatedAt: now, UpdatedAt: now})
	if err == nil {
		s.notify(ctx, o, "ORDER_CREATED")
	}
	return o, err
}
func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	return s.store.Get(ctx, id)
}
func (s *OrderService) List(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error) {
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
func (s *OrderService) Update(ctx context.Context, id uuid.UUID, name, email string, items []domain.OrderItem, status domain.OrderStatus) (domain.Order, error) {
	cleanName, cleanEmail, cleanItems, total, err := validate(name, email, items, status)
	if err != nil {
		return domain.Order{}, err
	}
	if status == "" {
		status = domain.StatusPending
	}
	o, err := s.store.Update(ctx, domain.Order{ID: id, CustomerName: cleanName, CustomerEmail: cleanEmail, Items: cleanItems, TotalAmount: total, Status: status, UpdatedAt: time.Now().UTC()})
	if err == nil {
		s.notify(ctx, o, "ORDER_UPDATED")
	}
	return o, err
}
func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.store.Delete(ctx, id)
	if err == nil {
		s.notify(ctx, domain.Order{ID: id}, "ORDER_DELETED")
	}
	return err
}
func (s *OrderService) notify(ctx context.Context, o domain.Order, event string) {
	if s.notifier == nil {
		return
	}
	// Notification is deliberately decoupled from the request deadline. The order
	// transaction has already committed, so a short independent timeout preserves
	// the best-effort semantics without making the REST mutation fail on outage.
	notifyCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
	defer cancel()
	if err := s.notifier.Send(notifyCtx, o.ID, event, fmt.Sprintf("Order %s %s", o.ID, strings.ToLower(strings.TrimPrefix(event, "ORDER_")))); err != nil {
		fmt.Printf("notification best-effort failure order=%s event=%s: %v\n", o.ID, event, err)
	}
}

func validate(name, email string, items []domain.OrderItem, status domain.OrderStatus) (string, string, []domain.OrderItem, float64, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" || len([]rune(name)) > 120 || email == "" {
		return "", "", nil, 0, ErrInvalidInput
	}
	if parsed, err := mail.ParseAddress(email); err != nil || parsed.Address != email {
		return "", "", nil, 0, ErrInvalidInput
	}
	if status != "" && status != domain.StatusPending && status != domain.StatusConfirmed && status != domain.StatusCancelled {
		return "", "", nil, 0, ErrInvalidInput
	}
	if len(items) == 0 {
		return "", "", nil, 0, ErrInvalidInput
	}
	clean := make([]domain.OrderItem, len(items))
	var total float64
	for i, item := range items {
		item.Name = strings.TrimSpace(item.Name)
		if item.Name == "" || item.Quantity <= 0 || item.UnitPrice < 0 {
			return "", "", nil, 0, ErrInvalidInput
		}
		clean[i] = item
		total += float64(item.Quantity) * item.UnitPrice
	}
	return name, email, clean, total, nil
}
