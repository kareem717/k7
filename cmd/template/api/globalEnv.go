package template

import (
	_ "embed"
)

//go:embed shared/files/globalenv.tmpl
var globalEnvTemplate []byte

func GlobalEnvTemplate() []byte {
	return globalEnvTemplate
}
