// Package program provides the
// main functionality of K7
package program

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	"github.com/kareem717/k7/cmd/program/api"
	"github.com/kareem717/k7/cmd/program/shared"
	"github.com/kareem717/k7/cmd/steps"
	"github.com/kareem717/k7/cmd/ui/multiinput"
	"github.com/kareem717/k7/cmd/ui/textinput"
	"github.com/spf13/cobra"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type Project struct {
	AbsolutePath string
	AppType      flags.App
	UnixBased    bool
}

// CreateAPIApp creates a new API app
func (p *Project) CreateAPIApp(name string, opts ...api.OptFunc) error {
	steps := steps.APISteps()
	appOpts := api.Options{
		AbsolutePath: ".",
		Framework:    apiFlags.Huma,
		DBMS:         apiFlags.Postgres,
		Git:          flags.Skip,
		UnixBased:    false,
	}

	for _, opt := range opts {
		if err := opt(&appOpts); err != nil {
			return fmt.Errorf("error applying option: %w", err)
		}
	}

	var shouldExit bool
	apiName := &textinput.Output{Output: name}
	tprogram := tea.NewProgram(textinput.InitialTextInputModel(apiName, "What is the name of your project?", &shouldExit))
	if _, err := tprogram.Run(); err != nil {
		log.Printf("Name of project contains an error: %v", err)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}

	projDir := filepath.Join(appOpts.AbsolutePath, apiName.Output)
	if shared.IsNonEmptyDir(projDir) {
		err := fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", apiName.Output)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	shared.Exit(tprogram, &shouldExit)

	// CREATE APP TYPE STEP
	for name, step := range steps.Steps {
		selection := &multiinput.Selection{}
		tprogram = tea.NewProgram(multiinput.InitialModelMulti(step.Options, selection, step.Headers, &shouldExit))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		// Retrieve the struct, modify it, and reassign it back to the map
		step.Field = selection.Choice
		steps.Steps[name] = step
		log.Printf("step.%s: %s", name, step.Field) // Logs the current step's field value
		shared.Exit(tprogram, &shouldExit)
	}

	log.Printf("steps: %+v", steps.Steps) // Logs the final state of all steps

	app, err := api.NewAPIApp(apiName.Output, appOpts, opts...)
	if err != nil {
		return fmt.Errorf("error creating API app: %w", err)
	}

	return app.CreateMainFile()
}
