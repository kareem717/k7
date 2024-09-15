/*
Copyright © 2024 Kareem Yakubu <kareem.717@icloud.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	"github.com/kareem717/k7/cmd/program"
	"github.com/kareem717/k7/cmd/program/api"
	"github.com/kareem717/k7/cmd/program/shared"
	stepsPkg "github.com/kareem717/k7/cmd/steps"
	"github.com/kareem717/k7/cmd/ui/multiinput"
	"github.com/kareem717/k7/cmd/ui/textinput"
	"github.com/spf13/cobra"
)

const logo = `
██ ▄█▀     ██████ ▓█████ ██▒   █▓▓█████  ███▄    █ 
██▄█▒    ▒██    ▒ ▓█   ▀▓██░   █▒▓█   ▀  ██ ▀█   █ 
▓███▄░    ░ ▓██▄   ▒███   ▓██  █▒░▒███   ▓██  ▀█ ██▒
▓██ █▄      ▒   ██▒▒▓█  ▄  ▒██ █░░▒▓█  ▄ ▓██▒  ▐▌██▒
▒██▒ █▄   ▒██████▒▒░▒████▒  ▒▀█░  ░▒████▒▒██░   ▓██░
▒ ▒▒ ▓▒   ▒ ▒▓▒ ▒ ░░░ ▒░ ░  ░ ▐░  ░░ ▒░ ░░ ▒░   ▒ ▒ 
░ ░▒ ▒░   ░ ░▒  ░ ░ ░ ░  ░  ░ ░░   ░ ░  ░░ ░░   ░ ▒░
░ ░░ ░    ░  ░  ░     ░       ░░     ░      ░   ░ ░ 
░  ░            ░     ░  ░     ░     ░  ░         ░ 
`

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ED6605")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
)

func init() {
	rootCmd.AddCommand(initCmd)
}

type Options struct {
	AppType *multiinput.Selection
	APIName *textinput.Output
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Quickly initialize a new project",
	Long:  "K7 is a CLI tool that allows users to spin up projects with the corresponding structure seamlessly.",

	Run: func(cmd *cobra.Command, args []string) {
		var tprogram *tea.Program
		var err error

		options := Options{
			AppType: &multiinput.Selection{},
			APIName: &textinput.Output{},
		}

		project := &program.Project{
			AppType: flags.AppAPI,
		}

		steps := stepsPkg.InitSteps()
		fmt.Printf("%s\n", logoStyle.Render(logo))

		var shouldExit bool

		// CREATE APP TYPE STEP
		// isInteractive := true
		step := steps.Steps[stepsPkg.AppType.String()]
		tprogram = tea.NewProgram(multiinput.InitialModelMulti(step.Options, options.AppType, step.Headers, &shouldExit))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		shared.Exit(tprogram, &shouldExit)

		// this type casting is always safe since the user interface can
		// only pass strings that can be cast to a flags.Framework instance
		project.AppType = flags.App(strings.ToLower(options.AppType.Choice))

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.AbsolutePath = currentWorkingDir

		switch project.AppType {
		case flags.AppAPI:
			project.CreateAPIApp(
				project.AbsolutePath,
				api.WithAbsolutePath(project.AbsolutePath),
				api.WithFramework(apiFlags.Huma),
				api.WithDBMS(apiFlags.Postgres),
				api.WithGit(flags.Skip),
				api.WithUnixBased(project.UnixBased),
			)
		}

		fmt.Println(endingMsgStyle.Render("\nNext steps:"))
		fmt.Println(endingMsgStyle.Render(fmt.Sprintf("• cd into the newly created project with: `cd %s`\n", project.AbsolutePath)))
	},
}