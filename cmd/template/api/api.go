package api

import (
	_ "embed"

	"github.com/kareem717/k7/cmd/template/shared"
)

// var (
// 	initPackages = []string{"go.uber.org/zap", "github.com/joho/godotenv"}

// 	pgxPostgresDriver = []string{"github.com/jackc/pgx/v5"}

// 	bunPackages         = []string{"github.com/uptrace/bun", "github.com/alexlast/bunzap"}
// 	bunPgDialectPackage = []string{"github.com/uptrace/bun/dialect/pgdialect"}

// 	humaPackage = []string{"github.com/danielgtaylor/huma/v2", "github.com/go-chi/chi/v5"}

// 	gotruePackage   = []string{"github.com/supabase-community/gotrue-go"}
// 	supabasePackage = []string{"github.com/supabase-community/supabase-go"}
// )

type TemplateOptions struct {
	storage StorageIndentifier
	server  ServerIndentifier
}

// Template represents a API app template
type Template struct {
	TemplateOptions
	main       []shared.Injectable
	extensions shared.GlobalFiles
}

type TemplateOption func(options *TemplateOptions)

func WithStorage(storage StorageIndentifier) TemplateOption {
	return func(options *TemplateOptions) {
		options.storage = storage
	}
}

func WithServer(server ServerIndentifier) TemplateOption {
	return func(options *TemplateOptions) {
		options.server = server
	}
}

func NewTemplate(opts ...TemplateOption) *Template {
	options := &TemplateOptions{
		storage: Postgres,
		server:  Huma,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &Template{
		TemplateOptions: *options,
	}
}

// Generate allows for generating an array of injectable templates
// for a customized API app template
func (t *Template) Generate() shared.AppTemplate {

	storageGenerator := NewStorageGenerator()
	storageInjectables := storageGenerator.Generate(t.storage)

	serverGenerator := NewServerGenerator()
	serverInjectables := serverGenerator.Generate(t.server)

	entitiesGenerator := NewEntitiesGenerator()
	entitiesInjectables := entitiesGenerator.Generate()

	serviceGenerator := NewServiceGenerator()
	serviceInjectables := serviceGenerator.Generate()

	allInjectables := append(storageInjectables, append(serverInjectables, append(entitiesInjectables, serviceInjectables...)...)...)

	return shared.AppTemplate{
		Templates: allInjectables,
	}
}

// Dependencies returns the dependencies for the API template
func (t *Template) Dependencies() []string {
	//TODO: implement
	return []string{}
}
