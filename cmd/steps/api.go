// Package steps provides utility for creating
// each step of the CLI
package steps

import flags "github.com/kareem717/k7/cmd/flags/api"

const (
	APIName      StepName = "Name"
	APIFramework StepName = "Framework"
	DBMS         StepName = "Database Management System"
	DBDriver     StepName = "Database Driver"
)

// InitSteps initializes and returns the *Steps to be used in the CLI program
func APISteps() *Steps {
	steps := &Steps{
		map[string]StepSchema{
			GitRepo.String(): {
				StepName: string(GitRepo),
				Headers:  "Which git option would you like to select for your project?",
				Options: []Item{
					{
						Title: "Commit",
						Desc:  "Initialize a new git repository and commit all the changes",
					},
					{
						Title: "Stage",
						Desc:  "Initialize a new git repository but only stage the changes",
					},
					{
						Title: "Skip",
						Desc:  "Proceed without initializing a git repository",
					},
				},
			},
			APIFramework.String(): {
				StepName: string(APIFramework),
				Options: []Item{
					{
						Title: "Huma",
						Desc:  "A modern, simple, fast & flexible micro framework for building HTTP REST/RPC APIs in Go backed by OpenAPI 3 and JSON Schema",
						Flag:  string(flags.Huma),
					},
				},
				Headers: "What API framework do you want to use in your Go project?",
			},
			DBMS.String(): {
				StepName: string(DBMS),
				Options: []Item{
					{
						Title: "Postgres",
						Desc:  "PostgreSQL is a powerful, open source object-relational database system",
						Flag:  string(flags.Postgres),
					},
				},
				Headers: "What database do you want to use in your Go project?",
			},
		},
	}

	return steps
}
