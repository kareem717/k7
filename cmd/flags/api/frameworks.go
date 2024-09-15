package flags

import (
	"fmt"
	"strings"
)

type Framework string

const (
	Huma Framework = "huma"
)

var AllowedFrameworks = []string{string(Huma)}

func (f Framework) String() string {
	return string(f)
}

func (f *Framework) Type() string {
	return "Framework"
}

func (f *Framework) Set(value string) error {
	for _, framework := range AllowedFrameworks {
		if framework == value {
			*f = Framework(value)
			return nil
		}
	}

	return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedFrameworks, ", "))
}
