/*
Copyright © 2024 Kareem Yakubu <kareem.717@icloud.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
	"github.com/kareem717/k7/cmd/program"
	stepsPkg "github.com/kareem717/k7/cmd/steps"
	"github.com/kareem717/k7/cmd/ui/multiInput"
	"github.com/kareem717/k7/cmd/ui/spinner"
	"github.com/kareem717/k7/cmd/ui/textinput"
	"github.com/kareem717/k7/cmd/utils"
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
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

func init() {
	rootCmd.AddCommand(initCmd)
}

type Options struct {
	AppType *multiInput.Selection
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
			AppType: &multiInput.Selection{},
			APIName: &textinput.Output{},
		}

		project := &program.Project{
			AppTypes: []flags.App{},
		}

		steps := stepsPkg.InitSteps()
		fmt.Printf("%s\n", logoStyle.Render(logo))

		// CREATE APP TYPE STEP
		isInteractive := true
		step := steps.Steps[stepsPkg.AppType.String()]
		tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, options.AppType, step.Headers, project))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.ExitCLI(tprogram)

		// this type casting is always safe since the user interface can
		// only pass strings that can be cast to a flags.Framework instance
		project.AppTypes = []flags.App{flags.App(strings.ToLower(options.AppType.Choice))}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.AbsolutePath = currentWorkingDir

		// This calls the templates
		for _, appType := range project.AppTypes {
			log.Printf("appType: %v", appType)
			switch appType {
			case flags.AppAPI:
				steps := stepsPkg.APISteps()
				isInteractive = true

				tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.APIName, "What is the name of your project?", project))
				if _, err := tprogram.Run(); err != nil {
					log.Printf("Name of project contains an error: %v", err)
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
				if doesDirectoryExistAndIsNotEmpty(options.APIName.Output) {
					err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", options.APIName.Output)
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
				project.ExitCLI(tprogram)

				name := options.APIName.Output
				if err != nil {
					log.Fatal("failed to set the name flag value", err)
				}

				// CREATE APP TYPE STEP
				for name, step := range steps.Steps {
					selection := &multiInput.Selection{}
					tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, selection, step.Headers, project))
					if _, err := tprogram.Run(); err != nil {
						cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
					}

					// Retrieve the struct, modify it, and reassign it back to the map
					step.Field = selection.Choice
					steps.Steps[name] = step
					log.Printf("step.%s: %s", name, step.Field) // Logs the current step's field value
					project.ExitCLI(tprogram)
				}

				log.Printf("steps: %+v", steps.Steps) // Logs the final state of all steps

				apiApp := project.CreateAPIApp(
					name,
					project.AbsolutePath,
					apiFlags.Framework(steps.Steps[stepsPkg.APIFramework.String()].Field),
					apiFlags.DBMS(steps.Steps[stepsPkg.DBMS.String()].Field),
					flags.Git(steps.Steps[stepsPkg.GitRepo.String()].Field),
					project.UnixBased,
				)

				log.Printf("apiApp: %+v", apiApp)

				err = apiApp.CreateMainFile()
			}
		}

		spinner := tea.NewProgram(spinner.InitialModelNew())

		// add synchronization to wait for spinner to finish
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("The program encountered an unexpected issue and had to exit. The error was:", r)
				fmt.Println("If you continue to experience this issue, please post a message on our GitHub page or join our Discord server for support.")
				if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
					log.Printf("Problem releasing terminal: %v", releaseErr)
				}
			}
		}()

		if err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				log.Printf("Problem releasing terminal: %v", releaseErr)
			}
			log.Printf("Problem creating files for project. %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		fmt.Println(endingMsgStyle.Render("\nNext steps:"))
		// TODO: change this to the actual project name
		fmt.Println(endingMsgStyle.Render(fmt.Sprintf("• cd into the newly created project with: `cd %s`\n", "idk")))

		if isInteractive {
			nonInteractiveCommand := utils.NonInteractiveCommand(cmd.Use, cmd.Flags())
			fmt.Println(tipMsgStyle.Render("Tip: Repeat the equivalent Blueprint with the following non-interactive command:"))
			fmt.Println(tipMsgStyle.Italic(false).Render(fmt.Sprintf("• %s\n", nonInteractiveCommand)))
		}
		err = spinner.ReleaseTerminal()
		if err != nil {
			log.Printf("Could not release terminal: %v", err)
			cobra.CheckErr(err)
		}
	},
}

// doesDirectoryExistAndIsNotEmpty checks if the directory exists and is not empty
func doesDirectoryExistAndIsNotEmpty(name string) bool {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Printf("could not read directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err))
		}
		if len(dirEntries) > 0 {
			return true
		}
	}
	return false
}
