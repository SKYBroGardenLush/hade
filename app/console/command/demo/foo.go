package demo

import (
	"github.com/SKYBroGardenLush/skycraper/framework/cobra"
	"log"
)

var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo的简要说明",
	Long:    "fool的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的列子",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}
