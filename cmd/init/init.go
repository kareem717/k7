package init

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	appType AppType
	appName string
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `Testing init command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Creating %s app\n", appType)

		if appType == "" {
			prompt := promptui.Prompt{
				Label: "App Name",
				Validate: func(input string) error {
					if len(input) < 3 {
						return errors.New("app name must be at least 3 characters")
					}
					return nil
				},
				Default: "my-app",
			}

			result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			appName = result
		}

		if appType == "" {
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U000029BF {{ .Name | cyan }}",
				Inactive: "  {{ .Name | cyan }}",
				Selected: "\U00002705 {{ .Name | red | cyan }}",
			}

			options := appType.Options()

			searcher := func(input string, index int) bool {
				option := options[index]
				name := strings.Replace(strings.ToLower(option.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			prompt := promptui.Select{
				Label:     "App Type",
				Items:     options,
				Templates: templates,
				Size:      4,
				Searcher:  searcher,
			}

			i, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			appType = options[i].Value

			fmt.Printf("You choose number %d: %s\n", i+1, appType)
		}

		fmt.Printf("Creating %s app...\n", appType)
	},
}

func init() {
	InitCmd.Flags().StringVarP(&appName, "app-name", "n", "", "The name of the application to create.")
	InitCmd.Flags().VarP(&appType, "app-type", "t", "The type of the application to create.")
}
