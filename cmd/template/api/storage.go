package api

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/template/shared"
)

var (
	//go:embed files/storage/*
	storageTemplateAssets embed.FS

	//go:embed files/storage/storage.go.tmpl
	storageInterfaceTemplate []byte
)

const (
	storageTemplateAssetDirectory = "files/storage"
	storageInterfaceTemplateName  = "storage.go"
)

type StorageIndentifier string

const (
	Postgres StorageIndentifier = "postgres"
)

type StorageGenerator struct {
	basePath string
}

func NewStorageGenerator() *StorageGenerator {
	//TODO: get from viper
	basePath := "internal/storage"

	return &StorageGenerator{basePath: basePath}
}

// Generate generates the storage injectables for a given storage indentifier
func (s *StorageGenerator) Generate(storage StorageIndentifier) []shared.Injectable {
	templates := []shared.Injectable{
		{
			FilePath: s.SubPath(storageInterfaceTemplateName),
			Bytes:    storageInterfaceTemplate,
		},
	}

	// Read through EVERY file in the storage directory recursively
	var readFiles func(dir string) error
	readFiles = func(dir string) error {
		entries, err := storageTemplateAssets.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())
			//TODO: check if storage is SELECTED storage implementatiion	
			if entry.IsDir() {
				if err := readFiles(fullPath); err != nil {
					return err
				}
			} else {
				byteData, err := storageTemplateAssets.ReadFile(fullPath)
				if err != nil {
					log.Fatalf("error reading file %s: %v", fullPath, err)
				}

				outputName := strings.Replace(fullPath[len(storageTemplateAssetDirectory)+1:], ".tmpl", "", 1)

				templates = append(templates, shared.Injectable{
					FilePath: s.SubPath(outputName),
					Bytes:    byteData,
				})
			}
		}
		return nil
	}

	if err := readFiles(filepath.Join(storageTemplateAssetDirectory, string(storage))); err != nil {
		log.Fatalf("error reading storage directory %s: %v", storageTemplateAssetDirectory, err)
	}

	return templates
}

// SubPath generates a path relative to the base path of the storage
func (s *StorageGenerator) SubPath(paths ...string) string {
	return filepath.Join(append([]string{s.basePath}, paths...)...)
}
