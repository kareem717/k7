package flags

import (
	"fmt"
	"strings"
)

type App string

const (
	AppAPI    App = "api"
)

var AllowedApps = []string{string(AppAPI)}

func (f App) String() string {
	return string(f)
}

func (f *App) Type() string {
	return "App"
}

func (f *App) Set(value string) error {
	for _, allowedApp := range AllowedApps {
		if allowedApp == value {
			*f = App(value)
			return nil
		}
	}

	return fmt.Errorf("app app to use. Allowed values: %s", strings.Join(AllowedApps, ", "))
}
