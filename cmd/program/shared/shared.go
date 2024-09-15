package shared

import (
	"log"
	"os"
	"runtime"

	"github.com/charmbracelet/bubbletea"
)

// IsUnixBased checks whether the current OS is Unix-based
func IsUnixBased() bool {
	if runtime.GOOS != "windows" {
		return true
	}

	return false
}

// Exit exits the CLI
func Exit(prog *tea.Program) {
	// logo render here
	if err := prog.ReleaseTerminal(); err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}
