package dbms

import (
	_ "embed"
)

//go:embed files/storage/storage.go.tmpl
var interfaceTemplate []byte

// InterfaceTemplate returns a byte slice that represents
// the storage (DBMS) interface.
func InterfaceTemplate() []byte {
	return interfaceTemplate
}

type Templater interface {
	Env() []byte
	Implementation() []byte
	InitialMigration() []byte
	Shared() []byte
	Foo() []byte
}
