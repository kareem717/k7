package dbms

import (
	_ "embed"
)

//go:embed files/storage/shared/shared.go.tmpl
var sharedFileTemplate []byte

//go:embed files/storage/storage.go.tmpl
var interfaceTemplate []byte

// SharedFileTemplate returns a byte slice that represents
// the file shared across all storage (DBMS) implementations.
func SharedFileTemplate() []byte {
	return sharedFileTemplate
}

// InterfaceTemplate returns a byte slice that represents
// the storage (DBMS) interface.
func InterfaceTemplate() []byte {
	return interfaceTemplate
}

type Templater interface {
	Implementation() []byte
	Env() []byte
	InitialMigration() []byte
}
