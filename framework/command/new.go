package command

import (
  "fmt"
  "github.com/SKYBroGardenLush/skycraper/framework/cobra"
  "github.com/SKYBroGardenLush/skycraper/framework/survey"
  "github.com/SKYBroGardenLush/skycraper/framework/utils"
  "os"
  "path/filepath"
)

var name string
var folder string
var mod string
var version string
var currentPath string

var newCommand = &cobra.Command{
  Use:   "new",
  Short: "脚手架命令",
  RunE: func(c *cobra.Command, args []string) error {

    //目录名称
    {
      prompt := &survey.Input{
        Message: "请输入目录名称:",
      }
      err := survey.AskOne(prompt, &name)
      if err != nil {
        return err
      }

      folder = filepath.Join(currentPath, name)
      if utils.Exists(folder) {
        isForce := false
        prompt2 := &survey.Confirm{
          Message: "目录" + folder + "已经存在，是否重新重建?(确认后执行)",
        }
        err := survey.AskOne(prompt2, &isForce)
        if err != nil {
          return err
        }
        if isForce {
          if err := os.RemoveAll(folder); err != nil {
            return err
          } else {
            fmt.Println("目录已存在，创建目录失败")
            return nil
          }
        }
      }
    }

    //mod module 名称
    {
      prompt := &survey.Input{
        Message: "请输入模块名称(go.mod 中的module,默认为文件夹名称):",
      }
      err := survey.AskOne(prompt, &mod)
      if err != nil {
        return err
      }
      if mod == "" {
        mod = name
      }
    }

    // 版本号
    {
      prompt := &survey.Input{
        Message: "请输入版本名称(参考 https://github.com/SKYBroGardenLush/skycraper/release,默认为最新版本):",
      }
      err := survey.AskOne(prompt, &version)
      if err != nil {
        return err
      }
      if version != "" {
        //确认版本号是否正确
      } else if version == "" {

      }
    }

    return nil
  },
}
