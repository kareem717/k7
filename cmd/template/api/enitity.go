package api

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/template/shared"
)

var (
	//go:embed files/entities/*
	entitiesTemplateAssets embed.FS
)

const (
	entitiesTemplateAssetDirectory = "files/entities"
	entitiesInterfaceTemplateName  = "entities.go"
)

type EntitiesGenerator struct {
	basePath string
}

func NewEntitiesGenerator() *EntitiesGenerator {
	//TODO: get from viper
	basePath := "internal/entities"

	return &EntitiesGenerator{
		basePath: basePath,
	}
}

// Generate generates the storage injectables for a given storage indentifier
func (s *EntitiesGenerator) Generate() []shared.Injectable {
	templates := []shared.Injectable{}

	// Read through EVERY file in the entities directory recursively
	var readFiles func(dir string) error
	readFiles = func(dir string) error {
		entries, err := entitiesTemplateAssets.ReadDir(dir)
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
				byteData, err := entitiesTemplateAssets.ReadFile(fullPath)
				if err != nil {
					log.Fatalf("error reading file %s: %v", fullPath, err)
				}

				outputName := strings.Replace(fullPath[len(entitiesTemplateAssetDirectory)+1:], ".tmpl", "", 1)

				templates = append(templates, shared.Injectable{
					FilePath: s.SubPath(outputName),
					Bytes:    byteData,
				})
			}
		}
		return nil
	}

	if err := readFiles(entitiesTemplateAssetDirectory); err != nil {
		log.Fatalf("error reading entities directory %s: %v", entitiesTemplateAssetDirectory, err)
	}

	return templates
}

// SubPath generates a path relative to the base path of the storage
func (s *EntitiesGenerator) SubPath(paths ...string) string {
	return filepath.Join(append([]string{s.basePath}, paths...)...)
}
