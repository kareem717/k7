package flags

import (
	"fmt"
	"strings"
)

type Framework string

// These are all the current frameworks supported. If you want to add one, you
// can simply copy and past a line here. Do not forget to also add it into the
// AllowedProjectTypes slice too!
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
	// Contains isn't available in 1.20 yet
	// if AllowedProjectTypes.Contains(value) {
	for _, framework := range AllowedFrameworks {
		if framework == value {
			*f = Framework(value)
			return nil
		}
	}

	return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedFrameworks, ", "))
}
