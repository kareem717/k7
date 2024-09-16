package api

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/kareem717/k7/cmd/template/shared"
)

var (
	//go:embed files/server/*
	serverTemplateAssets embed.FS
	//go:embed files/main
	mainFileTemplateAssets embed.FS
	//go:embed files/globals/*
	globalTemplateAssets    embed.FS
	serverInterfaceTemplate []byte
)

const (
	serverTemplateAssetDirectory   = "files/server"
	mainFileTemplateAssetDirectory = "files/main"
	globalTemplateAssetDirectory   = "files/globals"
	serverInterfaceTemplateName    = "server.go"
)

type ServerIndentifier string

const (
	Huma ServerIndentifier = "huma"
)

type ServerGenerator struct {
	basePath string
}

func NewServerGenerator() *ServerGenerator {
	//TODO: get from viper
	basePath := "internal/server"

	return &ServerGenerator{basePath: basePath}
}

// Generate generates the storage injectables for a given storage indentifier
func (s *ServerGenerator) Generate(server ServerIndentifier) []shared.Injectable {
	templates := []shared.Injectable{
		{
			FilePath: s.SubPath(serverInterfaceTemplateName),
			Bytes:    serverInterfaceTemplate,
		},
	}

	// Read through EVERY file in the server directory recursively
	var readFiles func(dir string) error
	readFiles = func(dir string) error {
		entries, err := serverTemplateAssets.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())
			// TODO: check if server is SELECTED server implementation
			if entry.IsDir() {
				if err := readFiles(fullPath); err != nil {
					return err
				}
			} else {
				byteData, err := serverTemplateAssets.ReadFile(fullPath)
				if err != nil {
					log.Fatalf("error reading file %s: %v", fullPath, err)
				}

				outputName := strings.Replace(fullPath[len(serverTemplateAssetDirectory)+1:], ".tmpl", "", 1)

				templates = append(templates, shared.Injectable{
					FilePath: s.SubPath(strings.Replace(outputName, string(server), "", 1)),
					Bytes:    byteData,
				})
			}
		}
		return nil
	}

	entries, err := mainFileTemplateAssets.ReadDir(mainFileTemplateAssetDirectory)
	if err != nil {
		log.Fatalf("error reading file %s: %v", mainFileTemplateAssetDirectory, err)
	}

	for _, entry := range entries {
		if strings.Replace(entry.Name(), ".tmpl", "", 1) == string(server)+".go" {
			fullPath := filepath.Join(mainFileTemplateAssetDirectory, entry.Name())
			byteData, err := mainFileTemplateAssets.ReadFile(fullPath)
			if err != nil {
				log.Fatalf("error reading file %s: %v", fullPath, err)
			}

			templates = append(templates, shared.Injectable{
				FilePath: "main.go",
				Bytes:    byteData,
			})
		}
	}

	entries, err = globalTemplateAssets.ReadDir(globalTemplateAssetDirectory)
	if err != nil {
		log.Fatalf("error reading file %s: %v", globalTemplateAssetDirectory, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(globalTemplateAssetDirectory, entry.Name())
		byteData, err := globalTemplateAssets.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("error reading file %s: %v", fullPath, err)
		}

		outputName := strings.Replace(entry.Name(), ".tmpl", "", 1)
		templates = append(templates, shared.Injectable{
			FilePath: outputName,
			Bytes:    byteData,
		})
	}


	if err := readFiles(serverTemplateAssetDirectory); err != nil {
		log.Fatalf("error reading server directory %s: %v", serverTemplateAssetDirectory, err)
	}
	return templates
}

// SubPath generates a path relative to the base path of the storage
func (s *ServerGenerator) SubPath(paths ...string) string {
	return filepath.Join(append([]string{s.basePath}, paths...)...)
}
