package shared

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type fileInjectionMethod string

const (
	mainMethod             fileInjectionMethod = "main"
	serverMethod           fileInjectionMethod = "server"
	routesMethod           fileInjectionMethod = "routes"
	DBMSMethod             fileInjectionMethod = "dbms"
	integrationTestsMethod fileInjectionMethod = "integration-tests"
	testsMethod            fileInjectionMethod = "tests"
	envMethod              fileInjectionMethod = "env"
)

// TemplateInjector is a struct that contains the framework and dbms maps
// and is used to inject templates into files
type TemplateInjector struct {
	basePath string      // The root path reference to the injector
	params   interface{} // The parameters to pass to the template
}

func NewTemplateInjector(basePath string, params interface{}) (*TemplateInjector, error) {
	return &TemplateInjector{basePath: basePath, params: params}, nil
}

// CreateFileWithInjection creates the given file at the
// project path, and injects the appropriate template


func (ti *TemplateInjector) Inject(input ...Injectable) error {
	for _, i := range input {
		fullPath := filepath.Join(ti.basePath, i.FilePath)

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(fullPath), 0770); err != nil {
			return fmt.Errorf("error creating directory %s: %w", filepath.Dir(fullPath), err)
		}

		createdFile, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("error creating file %s: %w", fullPath, err)
		}

		defer createdFile.Close()

		createdTemplate := template.Must(template.New(i.FilePath).Parse(string(i.Bytes)))

		err = createdTemplate.Execute(createdFile, ti.params)
		if err != nil {
			return fmt.Errorf("error executing template for file %s: %w", fullPath, err)
		}
	}

	return nil
}
