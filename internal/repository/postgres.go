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

func NewPostgres(ctx context.Context, databaseURL string, maxConns int32) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	config.MaxConns = maxConns
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Postgres{pool: pool}, nil
}

func (r *Postgres) Close() { r.pool.Close() }

func (r *Postgres) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO products (id,name,description,price,stock,created_at,updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id,name,description,price,stock,created_at,updated_at`, p.ID, p.Name, p.Description, p.Price, p.Stock, p.CreatedAt, p.UpdatedAt).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}
	return p, nil
}

func (r *Postgres) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var p domain.Product
	err := r.pool.QueryRow(ctx, `SELECT id,name,description,price,stock,created_at,updated_at FROM products WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, service.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (r *Postgres) List(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, `SELECT id,name,description,price,stock,created_at,updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()
	products := make([]domain.Product, 0, pageSize)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list products rows: %w", err)
	}
	return products, total, nil
}

func (r *Postgres) Update(ctx context.Context, p domain.Product) (domain.Product, error) {
	err := r.pool.QueryRow(ctx, `UPDATE products SET name=$2,description=$3,price=$4,stock=$5,updated_at=$6 WHERE id=$1
RETURNING id,name,description,price,stock,created_at,updated_at`, p.ID, p.Name, p.Description, p.Price, p.Stock, p.UpdatedAt).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, service.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}
	return p, nil
}

func (r *Postgres) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}
