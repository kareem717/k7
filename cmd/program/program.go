// Package program provides the
// main functionality of K7
package program

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	"github.com/kareem717/k7/cmd/program/api"
)

// A Project contains the data for the project folder
// being created, and methods that help with that process
type Project struct {
	Exit         bool
	AbsolutePath string
	AppTypes     []flags.App
	UnixBased    bool
}

// ExitCLI checks if the Project has been exited, and closes
// out of the CLI if it has
func (p *Project) ExitCLI(tprogram *tea.Program) {
	if p.Exit {
		// logo render here
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

// CreateAPIApp creates a new API app
func (p *Project) CreateAPIApp(
	name string,
	absolutePath string,
	framework apiFlags.Framework,
	dbms apiFlags.DBMS,
	gitOptions flags.Git,
	unixBased bool,
) api.APIApp {
	app := api.NewAPIApp(name, absolutePath, framework, dbms, gitOptions, unixBased)
	return app
}
