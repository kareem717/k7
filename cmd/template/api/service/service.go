package service

import (
	_ "embed"
)

//go:embed files/service.go.tmpl
var serviceInterfaceTemplate []byte

//go:embed files/domain/service.go.tmpl
var serviceImplementationTemplate []byte

//go:embed files/domain/foo/foo.go.tmpl
var fooServiceTemplate []byte

//go:embed files/domain/health/health.go.tmpl
var healthServiceTemplate []byte

type Templater struct {
}

func (e *Templater) Interface() []byte {
	return serviceInterfaceTemplate
}

func (e *Templater) Implementation() []byte {
	return serviceImplementationTemplate
}

func (e *Templater) Foo() []byte {
	return fooServiceTemplate
}

func (e *Templater) Health() []byte {
	return healthServiceTemplate
}
