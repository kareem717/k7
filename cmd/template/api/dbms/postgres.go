package dbms

import (
	_ "embed"
)

type PostgresTemplate struct{}

//go:embed files/storage/postgres/storage.go.tmpl
var postgresStorageTemplate []byte

//go:embed files/storage/postgres/migrations/2000000000000_init.sql.tmpl
var postgresMigrationTemplate []byte

//go:embed files/env/postgres.env.tmpl
var postgresEnvTemplate []byte

//go:embed files/storage/postgres/shared/shared.go.tmpl
var sharedFileTemplate []byte

//go:embed files/storage/postgres/foo/foo.go.tmpl
var fooFileTemplate []byte

// Implementation returns a byte slice that represents
// the postgres storage implementation.
func (m PostgresTemplate) Implementation() []byte {
	return postgresStorageTemplate
}

// Shared returns a byte slice that represents
// the shared file for postgres.
func (m PostgresTemplate) Shared() []byte {
	return sharedFileTemplate
}

// Shared returns a byte slice that represents
// the shared file for postgres.
func (m PostgresTemplate) Foo() []byte {
	return fooFileTemplate
}

// Env returns a byte slice that represents
// the postgres environment variables.
func (m PostgresTemplate) Env() []byte {
	return postgresEnvTemplate
}

// InitialMigration returns a byte slice that represents
// the initial postgres migration, providing a foo table and helper functions.
func (m PostgresTemplate) InitialMigration() []byte {
	return postgresMigrationTemplate
}
