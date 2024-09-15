// Package api provides utility functions that
// help with the templating of created API files.
package framework

import (
	_ "embed"
)

//go:embed files/air.toml.tmpl
var airTomlTemplate []byte

//go:embed files/README.md.tmpl
var readmeTemplate []byte

//go:embed files/makefile.tmpl
var makeTemplate []byte

//go:embed files/gitignore.tmpl
var gitIgnoreTemplate []byte

// MakeTemplate returns a byte slice that represents
// the default Makefile template.
func MakeTemplate() []byte {
	return makeTemplate
}

func GitIgnoreTemplate() []byte {
	return gitIgnoreTemplate
}

func AirTomlTemplate() []byte {
	return airTomlTemplate
}

// ReadmeTemplate returns a byte slice that represents
// the default README.md file template.
func ReadmeTemplate() []byte {
	return readmeTemplate
}

type Handler struct {
	Name    string
	Handler []byte
	Routes  []byte
}

type HandlerHelper struct {
	Auth    []byte
	Schemas []byte
}

type Handlers struct {
	Handlers []Handler
	Shared   HandlerHelper
}

type Middleware struct {
	Auth   []byte
	Shared []byte
}

// A Templater has the methods that help build the files
// in the Project folder, and is specific to a Framework
type Templater interface {
	Main() []byte
	Server() []byte
	Router() []byte
	Handlers() Handlers
	Middleware() Middleware
}
