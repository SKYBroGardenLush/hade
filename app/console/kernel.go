package console

import (
	"github.com/SKYBroGardenLush/skyscraper/app/console/command/demo"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/cobra"
	"github.com/SKYBroGardenLush/skyscraper/framework/command"
	"time"
)

func RunCommand(container framework.Container) error {
	//root Command
	var rootCmd = &cobra.Command{
		//定义root命令关键字
		Use:   "hade",
		Short: "hade 命令",
		Long:  "hade 框架提供命令行工具，使用这个命令行工具能方便执行框架自带命令，也能方便编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpCmd()
			return cmd.Help()
		},
		//不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	//为rootcmd设置服务容器
	rootCmd.SetContainer(container)
	//绑定框架命令
	command.AddKernelCommands(rootCmd)
	//绑定业务的命令`
	AddAppCommand(rootCmd)
	//执行RootCommand
	return rootCmd.Execute()
}

func AddAppCommand(root *cobra.Command) {
	root.AddDistributedCronCommand("foo_func_for_test", "@every 5s", demo.FooCommand, 2*time.Second)
	root.AddCommand(demo.FooCommand)
}
