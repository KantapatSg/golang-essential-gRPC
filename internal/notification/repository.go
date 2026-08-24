package notification

import (
	"context"
	"github.com/KantapatSg/golang-essential-gRPC/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{ pool *pgxpool.Pool }

func NewPostgres(ctx context.Context, url string, max int32) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = max
	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return &Postgres{pool: p}, nil
}
func (r *Postgres) Close() { r.pool.Close() }
func (r *Postgres) Create(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO notifications (id,order_id,event_type,message,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING created_at`, n.ID, n.OrderID, n.EventType, n.Message, n.CreatedAt).Scan(&n.CreatedAt)
	return n, err
}
func (r *Postgres) List(ctx context.Context, id uuid.UUID, page, size int) ([]domain.Notification, int64, error) {
	var total int64
	var err error
	if id == uuid.Nil {
		err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications`).Scan(&total)
	} else {
		err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE order_id=$1`, id).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}
	var rows pgx.Rows
	if id == uuid.Nil {
		rows, err = r.pool.Query(ctx, `SELECT id,order_id,event_type,message,created_at FROM notifications ORDER BY created_at DESC LIMIT $1 OFFSET $2`, size, (page-1)*size)
	} else {
		rows, err = r.pool.Query(ctx, `SELECT id,order_id,event_type,message,created_at FROM notifications WHERE order_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, id, size, (page-1)*size)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []domain.Notification{}
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.OrderID, &n.EventType, &n.Message, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, n)
	}
	return out, total, rows.Err()
}
