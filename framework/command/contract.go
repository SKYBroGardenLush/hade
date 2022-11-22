package command

import (
  "fmt"
  "github.com/SKYBroGardenLush/skycraper/framework/cobra"
  "github.com/SKYBroGardenLush/skycraper/framework/contract"
)

var envCommand = &cobra.Command{
  Use:   "env",
  Short: "获取当前的App环境",
  Run: func(cmd *cobra.Command, args []string) {
    //获取env环境
    container := cmd.Root().GetContainer()
    envService := container.MustMake(contract.EnvKey).(contract.Env)
    //打印环境
    fmt.Println("environment:", envService.AppEnv())
	},
}
