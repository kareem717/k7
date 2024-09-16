package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/flags"
	tpl "github.com/kareem717/k7/cmd/template/api"

	"github.com/kareem717/k7/cmd/program/shared"
	tplShared "github.com/kareem717/k7/cmd/template/shared"
	"github.com/kareem717/k7/cmd/utils"
	"github.com/spf13/cobra"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type APIApp struct {
	Name string
	shared.Options
}

func NewAPIApp(name string, opts shared.Options, optFuncs ...shared.OptFunc) (*APIApp, error) {
	app := APIApp{
		Name:    name,
		Options: opts,
	}

	for _, opt := range optFuncs {
		if err := opt(&app.Options); err != nil {
			return nil, fmt.Errorf("error applying option: %w", err)
		}
	}

	return &app, nil
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
	if shared.IsUnixBased() {
		p.UnixBased = true
	}

	// Create go.mod
	err = utils.InitGoMod(p.Name, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		cobra.CheckErr(err)
	}

	template := tpl.NewTemplate()

	err = utils.GoGetPackage(p.Name, template.Dependencies())
	if err != nil {
		log.Printf("Error getting dependencies: %v", err)
		cobra.CheckErr(err)
		return err
	}

	injector, err := tplShared.NewTemplateInjector(projectPath, p)
	if err != nil {
		log.Printf("Error creating template injector: %v", err)
		cobra.CheckErr(err)
		return err
	}

	templateFiles := template.Generate()

	err = injector.Inject(templateFiles.Templates...)
	if err != nil {
		log.Printf("Error injecting templates: %v", err)
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
