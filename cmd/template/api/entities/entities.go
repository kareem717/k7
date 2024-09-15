package entities

import (
	_ "embed"
)

//go:embed files/foo/foo.go.tmpl
var fooEntityTemplate []byte

//go:embed files/shared/timestamps.go.tmpl
var timestampsEntityTemplate []byte

type Templater struct {
}

func (e *Templater) Foo() []byte {
	return fooEntityTemplate
}

func (e *Templater) Timestamps() []byte {
	return timestampsEntityTemplate
}
