package framework

import (
	_ "embed"
)

//go:embed files/huma/router.go.tmpl
var humaRouterTemplate []byte

//go:embed files/huma/server.go.tmpl
var humaServerTemplate []byte

//go:embed files/huma/main.go.tmpl
var humaMainTemplate []byte

//go:embed files/huma/middleware/auth.go.tmpl
var humaAuthMiddlewareTemplate []byte

//go:embed files/huma/middleware/shared.go.tmpl
var humaSharedMiddlewareTemplate []byte

//go:embed files/huma/handler/foo/handler.go.tmpl
var humaFooHandlerTemplate []byte

//go:embed files/huma/handler/foo/routes.go.tmpl
var humaFooHandlerRoutesTemplate []byte

//go:embed files/huma/handler/health/handler.go.tmpl
var humaHealthHandlerTemplate []byte

//go:embed files/huma/handler/health/routes.go.tmpl
var humaHealthHandlerRoutesTemplate []byte

//go:embed files/huma/handler/shared/auth.go.tmpl
var humaHandlerAuthHelperTemplate []byte

//go:embed files/huma/handler/shared/schema.go.tmpl
var humaHandlerSchemaHelperTemplate []byte

// HumaTemplates contains the methods used for building
// an app that uses [github.com/danielgtaylor/huma]
type HumaTemplate struct{}

func (f HumaTemplate) Main() []byte {
	return humaMainTemplate
}
func (f HumaTemplate) Server() []byte {
	return humaServerTemplate
}

func (f HumaTemplate) Router() []byte {
	return humaRouterTemplate
}

func (f HumaTemplate) Handlers() Handlers {
	return Handlers{
		Handlers: []Handler{
			{
				Name: "foo",
				Handler: humaFooHandlerTemplate,
				Routes: humaFooHandlerRoutesTemplate,
			},
			{
				Name: "health",
				Handler: humaHealthHandlerTemplate,
				Routes: humaHealthHandlerRoutesTemplate,
			},
		},
		Shared: HandlerHelper{
			Auth: humaHandlerAuthHelperTemplate,
			Schemas: humaHandlerSchemaHelperTemplate,
		},
	}
}

func (f HumaTemplate) Middleware() Middleware {
	return Middleware{
		Auth: humaAuthMiddlewareTemplate,
		Shared: humaSharedMiddlewareTemplate,
	}
}
