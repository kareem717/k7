package steps

type StepSchema struct {
	StepName string // The name of a given step
	Options  []Item // The slice of each option for a given step
	Headers  string // The title displayed at the top of a given step
	Field    string // The field that is used to store the value of the step
}

// Steps contains a slice of steps
type Steps struct {
	Steps map[string]StepSchema
}

// An Item contains the data for each option
// in a StepSchema.Options
type Item struct {
	// Flag is the flag that is used to declare the value of the step
	Flag string
	// Title is the title of the option, if there is no flag, the title is used as the step value
	Title string
	// Desc is the description of the option
	Desc string
}

type StepName string

func (s StepName) String() string {
	return string(s)
}

const (
	AppType StepName = "App Type"
	GitRepo StepName = "Git Repository"
)

// InitSteps initializes and returns the *Steps to be used in the CLI program
func InitSteps() *Steps {
	steps := &Steps{
		map[string]StepSchema{
			AppType.String(): {
				StepName: string(AppType),
				Options: []Item{
					{
						Title: "API",
						Desc:  "A simple API server",
					},
					{
						Title: "Web",
						Desc:  "A simple web app",
					},
					{
						Title: "Mobile",
						Desc:  "A simple mobile app",
					},
				},
				Headers: "What type of app do you want to create?",
			},
		},
	}

	return steps
}
