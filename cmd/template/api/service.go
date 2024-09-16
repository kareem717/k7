package api

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/template/shared"
)

var (
	//go:embed files/service/*
	serviceTemplateAssets embed.FS



	//go:embed files/service/service.go.tmpl
	serviceInterfaceTemplate []byte
)

const (
	serviceTemplateAssetDirectory = "files/service"
	serviceInterfaceTemplateName  = "service.go"
)

type ServiceGenerator struct {
	basePath string
}

func NewServiceGenerator() *ServiceGenerator {
	//TODO: get from viper
	basePath := "internal/service"

	return &ServiceGenerator{basePath: basePath}
}

// Generate generates the storage injectables for a given storage indentifier
func (s *ServiceGenerator) Generate() []shared.Injectable {
	templates := []shared.Injectable{
		{
			FilePath: s.SubPath(serviceInterfaceTemplateName),
			Bytes:    serviceInterfaceTemplate,
		},
	}

	// Read through EVERY file in the service directory recursively
	var readFiles func(dir string) error
	readFiles = func(dir string) error {
		entries, err := serviceTemplateAssets.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				if err := readFiles(fullPath); err != nil {
					return err
				}
			} else {
				byteData, err := serviceTemplateAssets.ReadFile(fullPath)
				if err != nil {
					log.Fatalf("error reading file %s: %v", fullPath, err)
				}

				outputName := strings.Replace(fullPath[len(serviceTemplateAssetDirectory)+1:], ".tmpl", "", 1)

				templates = append(templates, shared.Injectable{
					FilePath: s.SubPath(outputName),
					Bytes:    byteData,
				})
			}
		}
		return nil
	}



	if err := readFiles(serviceTemplateAssetDirectory); err != nil {
		log.Fatalf("error reading service directory %s: %v", serviceTemplateAssetDirectory, err)
	}

	return templates
}

// SubPath generates a path relative to the base path of the storage
func (s *ServiceGenerator) SubPath(paths ...string) string {
	return filepath.Join(append([]string{s.basePath}, paths...)...)
}
