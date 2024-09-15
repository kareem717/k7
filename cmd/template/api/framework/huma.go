package framework

import (
	_ "embed"
)

//go:embed files/routes/huma.go.tmpl
var humaRoutesTemplate []byte

//go:embed files/server/huma.go.tmpl
var humaServerTemplate []byte

//go:embed files/main/huma_main.go.tmpl
var humaMainTemplate []byte

//go:embed files/tests/huma-test.go.tmpl
var humaTestHandlerTemplate []byte

// HumaTemplates contains the methods used for building
// an app that uses [github.com/danielgtaylor/huma]
type HumaTemplate struct{}

func (f HumaTemplate) Main() []byte {
	return humaMainTemplate
}
func (f HumaTemplate) Server() []byte {
	return humaServerTemplate
}

func (f HumaTemplate) Routes() []byte {
	return humaRoutesTemplate
}

func (f HumaTemplate) TestHandler() []byte {
	return humaTestHandlerTemplate
}
