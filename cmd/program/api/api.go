package api

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	tpl "github.com/kareem717/k7/cmd/template/api"

	"github.com/kareem717/k7/cmd/program/shared"
	"github.com/kareem717/k7/cmd/template/api/dbms"
	"github.com/kareem717/k7/cmd/template/api/entities"
	"github.com/kareem717/k7/cmd/template/api/framework"
	"github.com/kareem717/k7/cmd/template/api/service"
	"github.com/kareem717/k7/cmd/utils"
	"github.com/spf13/cobra"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type APIApp struct {
	Name         string
	FrameworkMap map[apiFlags.Framework]Framework
	DBMSMap      map[apiFlags.DBMS]DBMS
	Service      Service
	Entities     Entities
	shared.Options
}

func NewAPIApp(name string, opts shared.Options, optFuncs ...shared.OptFunc) (*APIApp, error) {
	app := APIApp{
		Name:         name,
		Options:      opts,
		FrameworkMap: make(map[apiFlags.Framework]Framework),
		DBMSMap:      make(map[apiFlags.DBMS]DBMS),
	}

	for _, opt := range optFuncs {
		if err := opt(&app.Options); err != nil {
			return nil, fmt.Errorf("error applying option: %w", err)
		}
	}

	return &app, nil
}

// A Framework contains the name and templater for a
// given Framework
type Framework struct {
	packageName []string
	templater   framework.Templater
}

type DBMS struct {
	name             string
	initialMigration string
	packageName      []string
	templater        dbms.Templater
}

type Service struct {
	packageName []string
	templater   service.Templater
}

type Entities struct {
	packageName []string
	templater   entities.Templater
}

var (
	initPackages = []string{"go.uber.org/zap", "github.com/joho/godotenv"}

	pgxPostgresDriver = []string{"github.com/jackc/pgx/v5"}

	bunPackages         = []string{"github.com/uptrace/bun", "github.com/alexlast/bunzap"}
	bunPgDialectPackage = []string{"github.com/uptrace/bun/dialect/pgdialect"}

	humaPackage = []string{"github.com/danielgtaylor/huma/v2", "github.com/go-chi/chi/v5"}

	gotruePackage   = []string{"github.com/supabase-community/gotrue-go"}
	supabasePackage = []string{"github.com/supabase-community/supabase-go"}
)

const (
	root = "/"

	internalServerPath   = "internal/server"
	internalStoreagePath = "internal/storage"
	internalServicePath  = "internal/service"
	internalEntitiesPath = "internal/entities"
)

// CheckOs checks Operation system and generates MakeFile and `go build` command
// Based on Project.Unixbase
func (p *APIApp) SetUnixBased() {
	if shared.IsUnixBased() {
		p.UnixBased = true
	}
}

// createFrameWorkMap adds the current supported
// Frameworks into a Project's FrameworkMap
func (p *APIApp) createFrameworkMap() {
	p.FrameworkMap[apiFlags.Huma] = Framework{
		packageName: humaPackage,
		templater:   framework.HumaTemplate{},
	}
}

func (p *APIApp) createDBMSMap() {
	p.DBMSMap[apiFlags.Postgres] = DBMS{
		name:             "postgres",
		initialMigration: "2000000000000_init.sql",
		//TODO: clean this up
		packageName: append(append(pgxPostgresDriver, bunPackages...), bunPgDialectPackage...),
		templater:   dbms.PostgresTemplate{},
	}
}

func (p *APIApp) createService() {
	p.Service = Service{
		packageName: humaPackage,
		templater:   service.Templater{},
	}
}

func (p *APIApp) createEntities() {
	p.Entities = Entities{
		packageName: humaPackage,
		templater:   entities.Templater{},
	}
}

// CreateMainFile creates the project folders and files,
// and writes to them depending on the selected options
func (p *APIApp) CreateMainFile() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0o754); err != nil {
			log.Printf("Could not create directory: %v", err)
			return err
		}
	}

	// Check if user.email is set.
	emailSet, err := utils.CheckGitConfig("user.email")
	if err != nil {
		cobra.CheckErr(err)
	}

	if !emailSet && p.Git.String() != flags.Skip {
		fmt.Println("user.email is not set in git config.")
		fmt.Println("Please set up git config before trying again.")
		panic("\nGIT CONFIG ISSUE: user.email is not set in git config.\n")
	}

	p.Name = strings.TrimSpace(p.Name)

	// Create a new directory with the project name
	projectPath := filepath.Join(p.AbsolutePath, p.Name)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := os.MkdirAll(projectPath, 0o751)
		if err != nil {
			log.Printf("Error creating root project directory %v\n", err)
			return err
		}
	}

	// Define Operating system
	p.SetUnixBased()

	// Create the map for our program
	p.createFrameworkMap()

	// Install the correct package for the selected driver
	p.createDBMSMap()

	p.createService()

	p.createEntities()

	// Create go.mod
	err = utils.InitGoMod(p.Name, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		cobra.CheckErr(err)
	}

	// Install the correct package for the selected framework
	err = utils.GoGetPackage(projectPath, p.FrameworkMap[p.Framework].packageName)
	if err != nil {
		log.Printf("Could not install go dependency for the chosen framework %v\n", err)
		cobra.CheckErr(err)
	}

	// Install the correct package for the selected DBMS
	err = utils.GoGetPackage(projectPath, p.DBMSMap[p.DBMS].packageName)
	if err != nil {
		log.Printf("Could not install go dependency for chosen DBMS %v\n", err)
		cobra.CheckErr(err)
	}

	// Install the correct package for the selected entities
	err = utils.GoGetPackage(projectPath, p.Entities.packageName)
	if err != nil {
		log.Printf("Could not install go dependency for chosen entities %v\n", err)
		cobra.CheckErr(err)
	}

	// Install the correct package for the selected entities
	err = utils.GoGetPackage(projectPath, p.Entities.packageName)
	if err != nil {
		log.Printf("Could not install go dependency for chosen entities %v\n", err)
		cobra.CheckErr(err)
	}

	injector, err := shared.NewTemplateInjector(projectPath, p)
	if err != nil {
		log.Printf("Error creating template injector: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create the DBMS.go file
	err = p.injectDBMSFiles(injector)
	if err != nil {
		log.Printf("Error injecting storage files: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.injectFrameworkFiles(injector)
	if err != nil {
		log.Printf("Error injecting framework files: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.injectWorkspaceFiles(injector)
	if err != nil {
		log.Printf("Error injecting workspace files: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.injectServiceFiles(injector)
	if err != nil {
		log.Printf("Error injecting service files: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.injectEntitiesFiles(injector)
	if err != nil {
		log.Printf("Error injecting entities files: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.injectEnvFile(injector)
	if err != nil {
		log.Printf("Error injecting env file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.GoTidy(projectPath)
	if err != nil {
		log.Printf("Could not go tidy in new project %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		log.Printf("Could not gofmt in new project %v\n", err)
		cobra.CheckErr(err)
		return err
	}

	if p.Git != flags.Skip {
		// Initialize git repo
		err = utils.ExecuteCmd("git", []string{"init"}, projectPath)
		if err != nil {
			log.Printf("Error initializing git repo: %v", err)
			cobra.CheckErr(err)
			return err
		}

		// Git add files
		err = utils.ExecuteCmd("git", []string{"add", "."}, projectPath)
		if err != nil {
			log.Printf("Error adding files to git repo: %v", err)
			cobra.CheckErr(err)
			return err
		}

		if p.Git == flags.Commit {
			// Git commit files
			err = utils.ExecuteCmd("git", []string{"commit", "-m", "Initial commit"}, projectPath)
			if err != nil {
				log.Printf("Error committing files to git repo: %v", err)
				cobra.CheckErr(err)
				return err
			}
		}
	}

	return nil
}

func (p *APIApp) injectDBMSFiles(ti *shared.TemplateInjector) error {
	dbmsConfig := p.DBMSMap[p.DBMS]

	// Create implementation agnostic helper file
	filePath := filepath.Join(internalStoreagePath, dbmsConfig.name, "shared", "shared.go")
	err := ti.Inject(filePath, dbmsConfig.templater.Shared())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	// Create the storage.go interface file
	filePath = filepath.Join(internalStoreagePath, "storage.go")
	err = ti.Inject(filePath, dbms.InterfaceTemplate())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	// Create the storage.go interface file
	filePath = filepath.Join(internalStoreagePath, dbmsConfig.name, "foo", "foo.go")
	err = ti.Inject(filePath, dbmsConfig.templater.Foo())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	// Create initial migration file
	initMigrationFile := fmt.Sprintf("%s/migrations/%s", dbmsConfig.name, dbmsConfig.initialMigration)
	filePath = filepath.Join(internalStoreagePath, initMigrationFile)
	err = ti.Inject(filePath, dbmsConfig.templater.InitialMigration())
	if err != nil {
		log.Printf("Error injecting migrations/0_foo.sql file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create DBMS specific implementation file
	implementationFile := fmt.Sprintf("%s/storage.go", dbmsConfig.name)
	filePath = filepath.Join(internalStoreagePath, implementationFile)
	err = ti.Inject(filePath, dbmsConfig.templater.Implementation())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *APIApp) injectWorkspaceFiles(ti *shared.TemplateInjector) error {
	// Create gitignore
	err := ti.Inject(".gitignore", framework.GitIgnoreTemplate())
	if err != nil {
		log.Printf("Error injecting .gitignore file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create air config file
	err = ti.Inject(".air.toml", framework.AirTomlTemplate())
	if err != nil {
		log.Printf("Error injecting .air.toml file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create Makefile
	err = ti.Inject("Makefile", framework.MakeTemplate())
	if err != nil {
		log.Printf("Error injecting Makefile file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create README.md
	err = ti.Inject("README.md", framework.ReadmeTemplate())
	if err != nil {
		log.Printf("Error injecting README.md file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *APIApp) injectFrameworkFiles(ti *shared.TemplateInjector) error {
	frameworkConfig := p.FrameworkMap[p.Framework]

	err := ti.Inject("main.go", frameworkConfig.templater.Main())
	if err != nil {
		log.Printf("Error injecting main.go file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	filePath := filepath.Join(internalServerPath, "server.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Server())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServerPath, "router.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Router())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServerPath, "middleware/auth.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Middleware().Auth)
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServerPath, "middleware/shared.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Middleware().Shared)
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	for _, handler := range frameworkConfig.templater.Handlers().Handlers {
		filePath = filepath.Join(internalServerPath, "handler", handler.Name, "handler.go")
		err = ti.Inject(filePath, handler.Handler)
		if err != nil {
			log.Printf("Error injecting %s file: %v", filePath, err)
			cobra.CheckErr(err)
			return err
		}

		filePath = filepath.Join(internalServerPath, "handler", handler.Name, "routes.go")
		err = ti.Inject(filePath, handler.Routes)
		if err != nil {
			log.Printf("Error injecting %s file: %v", filePath, err)
			cobra.CheckErr(err)
			return err
		}
	}

	filePath = filepath.Join(internalServerPath, "handler/shared/auth.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Handlers().Shared.Auth)
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServerPath, "handler/shared/schemas.go")
	err = ti.Inject(filePath, frameworkConfig.templater.Handlers().Shared.Schemas)
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *APIApp) injectServiceFiles(ti *shared.TemplateInjector) error {
	templater := p.Service.templater

	filePath := filepath.Join(internalServicePath, "service.go")
	err := ti.Inject(filePath, templater.Interface())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServicePath, "domain", "service.go")
	err = ti.Inject(filePath, templater.Implementation())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServicePath, "domain", "foo", "foo.go")
	err = ti.Inject(filePath, templater.Foo())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalServicePath, "domain", "health", "health.go")
	err = ti.Inject(filePath, templater.Health())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *APIApp) injectEntitiesFiles(ti *shared.TemplateInjector) error {
	templater := p.Entities.templater

	filePath := filepath.Join(internalEntitiesPath, "shared", "timestamps.go")
	err := ti.Inject(filePath, templater.Timestamps())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	filePath = filepath.Join(internalEntitiesPath, "foo", "foo.go")
	err = ti.Inject(filePath, templater.Foo())
	if err != nil {
		log.Printf("Error injecting %s file: %v", filePath, err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *APIApp) injectEnvFile(ti *shared.TemplateInjector) error {
	envBytes := [][]byte{
		tpl.GlobalEnvTemplate(),
		p.DBMSMap[p.DBMS].templater.Env(),
	}

	templ := bytes.Join(envBytes, []byte("\n"))

	err := ti.Inject(".env.example", templ)
	if err != nil {
		log.Printf("Error injecting .env.example file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = ti.Inject(".env.local", templ)
	if err != nil {
		log.Printf("Error injecting .env.local file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}
