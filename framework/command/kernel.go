package command

import "github.com/SKYBroGardenLush/skyscraper/framework/cobra"

// AddKernelCommands 将所有command/* 挂载到root command 中去
func AddKernelCommands(root *cobra.Command) {
	//挂载AppCommand 命令
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(envCommand)
	root.AddCommand(initBuildCommand())
	//挂载dev 调试命令
	root.AddCommand(initDevCommand())
	//挂载provider相关命令
	root.AddCommand(initProviderCommand())
	//挂载cmd相关命令
	root.AddCommand(initCmdCommand())
	//挂载中间件相关命令
	root.AddCommand(initMiddlewareCommand())
	//挂载swagger相关命令
	root.AddCommand(initSwaggerCommand())
}
