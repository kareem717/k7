package shared

import (
	"log"
	"os"
	"runtime"

	"github.com/charmbracelet/bubbletea"
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
