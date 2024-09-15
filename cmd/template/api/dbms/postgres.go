package dbms

import (
	_ "embed"
)

type PostgresTemplate struct{}

//go:embed files/storage/postgres/storage.go.tmpl
var postgresStorageTemplate []byte

//go:embed files/storage/postgres/migrations/20240828191000_create_foo_table.sql.tmpl
var postgresMigrationTemplate []byte

//go:embed files/env/postgres.env.tmpl
var postgresEnvTemplate []byte

// Implementation returns a byte slice that represents
// the postgres storage implementation.
func (m PostgresTemplate) Implementation() []byte {
	return postgresStorageTemplate
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

