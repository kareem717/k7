package api

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	tpl "github.com/kareem717/k7/cmd/template/api"

	"github.com/kareem717/k7/cmd/program/shared"
	"github.com/kareem717/k7/cmd/template/api/dbms"
	"github.com/kareem717/k7/cmd/template/api/framework"
	"github.com/kareem717/k7/cmd/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	Framework    apiFlags.Framework
	DBMS         apiFlags.DBMS
	Git          flags.Git
	UnixBased    bool
	AbsolutePath string
}

type OptFunc func(app *Options) error

func WithFramework(f apiFlags.Framework) OptFunc {
	return func(app *Options) error {
		app.Framework = f
		return nil
	}
}

func WithDBMS(d apiFlags.DBMS) OptFunc {
	return func(app *Options) error {
		app.DBMS = d
		return nil
	}
}

func WithGit(g flags.Git) OptFunc {
	return func(app *Options) error {
		app.Git = g
		return nil
	}
}

// WithUnixBased sets the UnixBased flag to true
func WithUnixBased(b bool) OptFunc {
	return func(app *Options) error {
		app.UnixBased = b
		return nil
	}
}

func WithAbsolutePath(path string) OptFunc {
	return func(app *Options) error {
		app.AbsolutePath = path
		return nil
	}
}

// A Project contains the data for the project folder
// being created, and methods that help with that process
type APIApp struct {
	Name         string
	FrameworkMap map[apiFlags.Framework]Framework
	DBMSMap      map[apiFlags.DBMS]DBMS
	Options
}

func NewAPIApp(name string, opts Options, optFuncs ...OptFunc) (*APIApp, error) {
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
	templater   Templater
}

type DBMS struct {
	packageName []string
	templater   DBMSTemplater
}

// A Templater has the methods that help build the files
// in the Project folder, and is specific to a Framework
type Templater interface {
	Main() []byte
	Server() []byte
	Routes() []byte
	TestHandler() []byte
}

type DBMSTemplater interface {
	Env() []byte
	Implementation() []byte
	InitialMigration() []byte
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
		//TODO: clean this up
		packageName: append(append(pgxPostgresDriver, bunPackages...), bunPgDialectPackage...),
		templater:   dbms.PostgresTemplate{},
	}
}

// CreateMainFile creates the project folders and files,
// and writes to them depending on the selected options
func (p *APIApp) CreateMainFile() error {
	log.Printf("\n\np: %+v\n\n", p)
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

	// Create the storage folder
	err = p.CreatePath(internalStoreagePath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", internalStoreagePath)
		cobra.CheckErr(err)
		return err
	}

	// Create the DBMS.go file
	err = p.CreateFileWithInjection(internalStoreagePath, projectPath, "storage.go", "DBMS")
	if err != nil {
		log.Printf("Error injecting DBMS.go file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(root, projectPath, "main.go", "main")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	makeFile, err := os.Create(filepath.Join(projectPath, "Makefile"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	defer makeFile.Close()

	// inject makefile template
	makeFileTemplate := template.Must(template.New("makefile").Parse(string(framework.MakeTemplate())))
	err = makeFileTemplate.Execute(makeFile, p)
	if err != nil {
		return err
	}

	readmeFile, err := os.Create(filepath.Join(projectPath, "README.md"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer readmeFile.Close()

	// inject readme template
	readmeFileTemplate := template.Must(template.New("readme").Parse(string(framework.ReadmeTemplate())))
	err = readmeFileTemplate.Execute(readmeFile, p)
	if err != nil {
		return err
	}

	err = p.CreatePath(internalServerPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", internalServerPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(internalServerPath, projectPath, "routes.go", "routes")
	if err != nil {
		log.Printf("Error injecting routes.go file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(internalServerPath, projectPath, "routes_test.go", "tests")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(internalServerPath, projectPath, "server.go", "server")
	if err != nil {
		log.Printf("Error injecting server.go file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(root, projectPath, ".env", "env")
	if err != nil {
		log.Printf("Error injecting .env file: %v", err)
		cobra.CheckErr(err)
		return err
	}

	// Create gitignore
	gitignoreFile, err := os.Create(filepath.Join(projectPath, ".gitignore"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer gitignoreFile.Close()

	// inject gitignore template
	gitignoreTemplate := template.Must(template.New(".gitignore").Parse(string(framework.GitIgnoreTemplate())))
	err = gitignoreTemplate.Execute(gitignoreFile, p)
	if err != nil {
		return err
	}

	// Create .air.toml file
	airTomlFile, err := os.Create(filepath.Join(projectPath, ".air.toml"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	defer airTomlFile.Close()

	// inject air.toml template
	airTomlTemplate := template.Must(template.New("airtoml").Parse(string(framework.AirTomlTemplate())))
	err = airTomlTemplate.Execute(airTomlFile, p)
	if err != nil {
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

	nameSet, err := utils.CheckGitConfig("user.name")
	if err != nil {
		cobra.CheckErr(err)
	}

	if p.Git.String() != flags.Skip {
		if !nameSet {
			fmt.Println("user.name is not set in git config.")
			fmt.Println("Please set up git config before trying again.")
			panic("\nGIT CONFIG ISSUE: user.name is not set in git config.\n")
		}
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

		if p.Git.String() == flags.Commit {
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

// CreatePath creates the given directory in the projectPath
func (p *APIApp) CreatePath(pathToCreate string, projectPath string) error {
	path := filepath.Join(projectPath, pathToCreate)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0o751)
		if err != nil {
			log.Printf("Error creating directory %v\n", err)
			return err
		}
	}

	return nil
}

// CreateFileWithInjection creates the given file at the
// project path, and injects the appropriate template
func (p *APIApp) CreateFileWithInjection(pathToCreate string, projectPath string, fileName string, methodName string) error {
	createdFile, err := os.Create(filepath.Join(projectPath, pathToCreate, fileName))
	if err != nil {
		return err
	}

	defer createdFile.Close()

	switch methodName {
	case "main":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.Framework].templater.Main())))
		err = createdTemplate.Execute(createdFile, p)
	case "server":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.Framework].templater.Server())))
		err = createdTemplate.Execute(createdFile, p)
	case "routes":
		routeFileBytes := p.FrameworkMap[p.Framework].templater.Routes()
		createdTemplate := template.Must(template.New(fileName).Parse(string(routeFileBytes)))
		err = createdTemplate.Execute(createdFile, p)
	case "DBMS":
		log.Printf("createdTemplate: %v", "there")
		log.Printf("templater: %v", p.DBMSMap[p.DBMS].templater)
		log.Printf("driver: %v", p.DBMS)

		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBMSMap[p.DBMS].templater.Implementation())))
		log.Printf("createdTemplate: %v", "here")
		err = createdTemplate.Execute(createdFile, p)
	case "integration-tests":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.DBMSMap[p.DBMS].templater.InitialMigration())))
		err = createdTemplate.Execute(createdFile, p)
	case "tests":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.Framework].templater.TestHandler())))
		err = createdTemplate.Execute(createdFile, p)
	case "env":
		if p.DBMS != "none" {

			envBytes := [][]byte{
				tpl.GlobalEnvTemplate(),
				p.DBMSMap[p.DBMS].templater.Env(),
			}
			createdTemplate := template.Must(template.New(fileName).Parse(string(bytes.Join(envBytes, []byte("\n")))))
			err = createdTemplate.Execute(createdFile, p)

		} else {
			createdTemplate := template.Must(template.New(fileName).Parse(string(tpl.GlobalEnvTemplate())))
			err = createdTemplate.Execute(createdFile, p)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
