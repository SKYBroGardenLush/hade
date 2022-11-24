package command

import (
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework/cobra"
	"log"
	"os/exec"
)

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译文件",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "编译hade命令",
	RunE: func(c *cobra.Command, args []string) error {
		//获取path路径下的npm命令
		path, err := exec.LookPath("go")
		if err != nil {
			log.Println("请安装go在你PATH路径下")
		}

		//=========适用windows系统============
		cmd := exec.Command(path, "build", "-o", "hade.exe", ".\\main.go") // 适用于windows系统
		//=========适用windows系统============

		//=========适用linux系统==============
		//cmd := exec.Command(path, "build", "-o", "hade", ".)
		//=========适用linux系统==============

		//将输出保存到out中
		out, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Println("===========go编译失败==============")
			fmt.Println(string(out))
			fmt.Println("=================================")
			return err
		}
		fmt.Print(string(out))
		fmt.Println("============go编译成功==============")
		fmt.Println("====run ./hade direct=============")
		return nil
	},
}

//编译后端
var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args)
	},
}

//打印前端的命令
var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "使用npm编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		//获取path路径下的npm命令
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Println("请安装npm在你PATH路径下")
		}

		//执行npm run build
		cmd := exec.Command(path, "run", "build")
		//将输出保存到out中
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("===========前端编译失败==============")
			fmt.Println(string(out))
			fmt.Println("===========前端编译失败==============")
			return err
		}
		fmt.Print(string(out))
		fmt.Println("============前端编译成功==============")
		return nil

	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildAllCommand)
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildBackendCommand)
	return buildCommand
}
