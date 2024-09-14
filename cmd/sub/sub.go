package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	test string
)

var SubCmd = &cobra.Command{
	Use:   "sub",
	Short: "A brief description of your command",
	Long:  `Testing sub command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sub called")
	},
}

func init() {
	SubCmd.Flags().StringVarP(&test, "test", "t", "", "test flag")

	if err := SubCmd.MarkFlagRequired("test"); err != nil {
		fmt.Println(err)
	}
}
