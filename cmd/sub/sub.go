package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var SubCmd = &cobra.Command{
	Use:   "sub",
	Short: "A brief description of your command",
	Long:  `Testing sub command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sub called")
	},
}

func init() {}
