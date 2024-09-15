package storage

import (
	"context"

	"kareem/internal/entities/foo"
	"kareem/internal/storage/postgres/shared"
)

type FooRepository interface {
	Create(ctx context.Context, params foo.CreateFooParams) (foo.Foo, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, paginationParams shared.PaginationRequest) ([]foo.Foo, error)
	Update(ctx context.Context, id int, params foo.UpdateFooParams) (foo.Foo, error)
	GetById(ctx context.Context, id int) (foo.Foo, error)
}

type RepositoryProvider interface {
	Foo() FooRepository
}

type Transaction interface {
	RepositoryProvider
	Commit() error
	Rollback() error
	SubTransaction() (Transaction, error)
}

type Repository interface {
	RepositoryProvider
	HealthCheck(ctx context.Context) error
	NewTransaction() (Transaction, error)
	RunInTx(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}
