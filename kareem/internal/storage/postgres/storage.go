package postgres

import (
	"context"
	"errors"
	"time"

	"kareem/internal/storage"
	"kareem/internal/storage/postgres/foo"

	"github.com/alexlast/bunzap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"
)

type Config struct {
	URL                   string
	MaxConnections        int32
	MinConnections        int32
	MaxConnectionIdleTime time.Duration
	MaxConnectionLifetime time.Duration
}

type ConfigOption func(*Config)

func NewConfig(url string, options ...ConfigOption) Config {
	config := Config{
		URL:                   url,
		MaxConnections:        10,
		MinConnections:        1,
		MaxConnectionIdleTime: 1 * time.Hour,
		MaxConnectionLifetime: 1 * time.Hour,
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

func WithMaxConnections(maxConnections int32) ConfigOption {
	return func(c *Config) {
		c.MaxConnections = maxConnections
	}
}

func WithMinConnections(minConnections int32) ConfigOption {
	return func(c *Config) {
		c.MinConnections = minConnections
	}
}

func WithMaxConnectionIdleTime(maxConnectionIdleTime time.Duration) ConfigOption {
	return func(c *Config) {
		c.MaxConnectionIdleTime = maxConnectionIdleTime
	}
}

func WithMaxConnectionLifetime(maxConnectionLifetime time.Duration) ConfigOption {
	return func(c *Config) {
		c.MaxConnectionLifetime = maxConnectionLifetime
	}
}

func configDBPool(config Config) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = config.MaxConnections
	poolConfig.MinConns = config.MinConnections
	poolConfig.MaxConnIdleTime = config.MaxConnectionIdleTime
	poolConfig.MaxConnLifetime = config.MaxConnectionLifetime

	return poolConfig, nil
}

type transaction struct {
	fooRepo *foo.FooRepository
	tx      *bun.Tx
	ctx     context.Context
}

func (t *transaction) Foo() storage.FooRepository {
	return t.fooRepo
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *transaction) SubTransaction() (storage.Transaction, error) {
	tx, err := t.tx.BeginTx(t.ctx, nil)
	if err != nil {
		return nil, err
	}

	return &transaction{
		fooRepo: foo.NewFooRepository(tx, t.ctx),
		tx:      &tx,
	}, nil
}

type Repository struct {
	fooRepo *foo.FooRepository
	db      *bun.DB
	ctx     context.Context
}

func NewRepository(config Config, ctx context.Context, logger *zap.Logger) *Repository {
	poolConfig, err := configDBPool(config)
	if err != nil {
		logger.Fatal("Error creating pool config: %v", zap.Error(err))
	}

	sqldb := stdlib.OpenDB(*poolConfig.ConnConfig)
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger:       logger,
		SlowDuration: 200 * time.Millisecond, // Omit to log all operations as debug
	}))

	// Increase timeout duration
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Attempting to ping the database...")
	err = db.PingContext(pingCtx)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			logger.Fatal("ping was canceled by the client: %v", zap.Error(err))
		case errors.Is(err, context.DeadlineExceeded):
			logger.Fatal("ping timed out: %v", zap.Error(err))
		default:
			logger.Fatal("ping failed: %v", zap.Error(err))
		}
	}

	logger.Info("Successfully connected to the database.")
	return &Repository{
		fooRepo: foo.NewFooRepository(db, ctx),
		db:      db,
		ctx:     ctx,
	}
}

func (r *Repository) Foo() storage.FooRepository {
	return r.fooRepo
}

func (r *Repository) HealthCheck(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) NewTransaction() (storage.Transaction, error) {
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return nil, err
	}

	return &transaction{
		fooRepo: foo.NewFooRepository(tx, r.ctx),
		tx:      &tx,
		ctx:     r.ctx,
	}, nil
}

func (r *Repository) RunInTx(ctx context.Context, fn func(ctx context.Context, tx storage.Transaction) error) error {
	tx, err := r.NewTransaction()
	if err != nil {
		return err
	}

	err = fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
