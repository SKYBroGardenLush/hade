package user

import (
	"github.com/SKYBroGardenLush/skycraper/framework/cobra"
)

var UserCommand = &cobra.Command{
	Use:     "user",
	Short:   "user的简要说明",
	Long:    "user的长说明",
	Example: "user命令的列子",
	RunE: func(cmd *cobra.Command, args []string) error {
		//书写你的逻辑
		return nil
	},
}
