package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/KantapatSg/golang-essential-gRPC/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{ pool *pgxpool.Pool }

func NewPostgres(ctx context.Context, url string, maxConns int32) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	cfg.MaxConns = maxConns
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}
func (r *Postgres) Close() { r.pool.Close() }
func (r *Postgres) Create(ctx context.Context, o domain.Order) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer tx.Rollback(ctx)
	if err = tx.QueryRow(ctx, `INSERT INTO orders (id,customer_name,customer_email,total_amount,status,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at,updated_at`, o.ID, o.CustomerName, o.CustomerEmail, o.TotalAmount, o.Status, o.CreatedAt, o.UpdatedAt).Scan(&o.CreatedAt, &o.UpdatedAt); err != nil {
		return domain.Order{}, fmt.Errorf("create order: %w", err)
	}
	if err = insertItems(ctx, tx, o.ID, o.Items); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return o, nil
}
func insertItems(ctx context.Context, tx pgx.Tx, id uuid.UUID, items []domain.OrderItem) error {
	for _, i := range items {
		if _, err := tx.Exec(ctx, `INSERT INTO order_items (order_id,name,quantity,unit_price) VALUES ($1,$2,$3,$4)`, id, i.Name, i.Quantity, i.UnitPrice); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}
	return nil
}
func (r *Postgres) Get(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	var o domain.Order
	err := r.pool.QueryRow(ctx, `SELECT id,customer_name,customer_email,total_amount,status,created_at,updated_at FROM orders WHERE id=$1`, id).Scan(&o.ID, &o.CustomerName, &o.CustomerEmail, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, service.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if err = r.items(ctx, &o); err != nil {
		return domain.Order{}, err
	}
	return o, nil
}
func (r *Postgres) items(ctx context.Context, o *domain.Order) error {
	rows, err := r.pool.Query(ctx, `SELECT name,quantity,unit_price FROM order_items WHERE order_id=$1 ORDER BY id`, o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	o.Items = []domain.OrderItem{}
	for rows.Next() {
		var i domain.OrderItem
		if err := rows.Scan(&i.Name, &i.Quantity, &i.UnitPrice); err != nil {
			return err
		}
		o.Items = append(o.Items, i)
	}
	return rows.Err()
}
func (r *Postgres) List(ctx context.Context, page, size int) ([]domain.Order, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT id,customer_name,customer_email,total_amount,status,created_at,updated_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	result := make([]domain.Order, 0, size)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.CustomerEmail, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if err := r.items(ctx, &o); err != nil {
			return nil, 0, err
		}
		result = append(result, o)
	}
	return result, total, rows.Err()
}
func (r *Postgres) Update(ctx context.Context, o domain.Order) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `UPDATE orders SET customer_name=$2,customer_email=$3,total_amount=$4,status=$5,updated_at=$6 WHERE id=$1 RETURNING created_at,updated_at`, o.ID, o.CustomerName, o.CustomerEmail, o.TotalAmount, o.Status, o.UpdatedAt).Scan(&o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, service.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id=$1`, o.ID); err != nil {
		return domain.Order{}, err
	}
	if err = insertItems(ctx, tx, o.ID, o.Items); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return o, nil
}
func (r *Postgres) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}
