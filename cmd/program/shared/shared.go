package shared

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/charmbracelet/bubbletea"
	"github.com/kareem717/k7/cmd/ui/spinner"
	"github.com/kareem717/k7/cmd/ui/textinput"
	"github.com/spf13/cobra"
)

// IsUnixBased checks whether the current OS is Unix-based
func IsUnixBased() bool {
	if runtime.GOOS != "windows" {
		return true
	}

	return false
}

// Exit exits the CLI
func Exit(prog *tea.Program, shouldExit *bool) {
	//TODO: handle this better
	if *shouldExit {
		// logo render here
		if err := prog.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

// IsNonEmptyDir checks if the directory exists and is not empty,
// if an error occurs a message is printed to the console and the program exits
func IsNonEmptyDir(name string) bool {
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

// CreatePath creates a directory within the given projectPath
func CreatePath(newDir string, projectPath string) error {
	path := filepath.Join(projectPath, newDir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0o751)
		if err != nil {
			log.Printf("Error creating directory %v\n", err)
			return err
		}
	}

	return nil
}

func RunWithSpinner(prog func() error) error {
	// TODO: fix: If the prog fails, the terminal gets messed up
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

	// This calls the templates
	err := prog()
	if err != nil {
		if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
			log.Printf("Problem releasing terminal: %v", releaseErr)
		}
		log.Printf("Problem creating files for project. %v", err)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}

	err = spinner.ReleaseTerminal()
	if err != nil {
		log.Printf("Could not release terminal: %v", err)
		cobra.CheckErr(err)
	}

	return nil
}
