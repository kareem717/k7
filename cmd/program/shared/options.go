package shared

import (
	"github.com/kareem717/k7/cmd/flags"
	apiFlags "github.com/kareem717/k7/cmd/flags/api"
)

type Options struct {
	Framework    apiFlags.Framework
	DBMS         apiFlags.DBMS
	Git          flags.Git
	UnixBased    bool
	AbsolutePath string
}

type OptFunc func(app *Options) error

func WithFramework(f apiFlags.Framework) OptFunc {
	return func(app *Options) error {
		app.Framework = f
		return nil
	}
}

func WithDBMS(d apiFlags.DBMS) OptFunc {
	return func(app *Options) error {
		app.DBMS = d
		return nil
	}
}

func WithGit(g flags.Git) OptFunc {
	return func(app *Options) error {
		app.Git = g
		return nil
	}
}

// WithUnixBased sets the UnixBased flag to true
func WithUnixBased(b bool) OptFunc {
	return func(app *Options) error {
		app.UnixBased = b
		return nil
	}
}

func WithAbsolutePath(path string) OptFunc {
	return func(app *Options) error {
		app.AbsolutePath = path
		return nil
	}
}
