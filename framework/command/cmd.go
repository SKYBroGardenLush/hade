package command

import (
  "fmt"
  "github.com/pkg/errors"
  "github.com/SKYBroGardenLush/skycraper/framework/cobra"
  "github.com/SKYBroGardenLush/skycraper/framework/contract"
  "github.com/SKYBroGardenLush/skycraper/framework/survey"
  "github.com/SKYBroGardenLush/skycraper/framework/utils"
  "golang.org/x/text/cases"
  "golang.org/x/text/language"
  "os"
  "path/filepath"
  "text/template"
)

func initCmdCommand() *cobra.Command {
  cmdCommand.AddCommand(cmdListCommand)
  cmdCommand.AddCommand(cmdNewCommand)
  return cmdCommand
}

var cmdCommand = &cobra.Command{
  Use:   "command",
  Short: "创键命令行",
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.Help()
    return nil
  },
}

var cmdListCommand = &cobra.Command{
  Use:   "list",
  Short: "列出所有控制台命令",
  RunE: func(cmd *cobra.Command, args []string) error {
    cmds := cmd.Root().Commands()
    ps := [][]string{}
    for _, cmd := range cmds {
      line := []string{cmd.Name(), cmd.Short}
      ps = append(ps, line)
    }
    for _, line := range ps {
      fmt.Printf("%-8s   %s\n", line[0], line[1])
    }
    return nil
  },
}

var cmdNewCommand = &cobra.Command{
  Use:     "new",
  Short:   "创键一个新的命令行",
  Aliases: []string{"create", "init"},
  RunE: func(c *cobra.Command, args []string) error {
    container := c.GetContainer()
    fmt.Println("创键一个服务")
    var name string
    //命令名称
    {
      prompt := &survey.Input{
        Message: "请输入新命令名称",
      }
      err := survey.AskOne(prompt, &name)
      if err != nil {
        return err
      }
    }

    //检查命令是否存在
    app := container.MustMake(contract.AppKey).(contract.App)
    appFolder := app.AppFolder()
    cmdFolder := filepath.Join(appFolder, "console", "command")
    subDirs, err := utils.SubDir(cmdFolder)

    if err != nil {
      return err
    }

    subColl := utils.NewStrCollection(subDirs)
    if subColl.Contains(name) {
      fmt.Println("命令名称已经存在")
      return nil
    }

    //创键命令文件夹
    if err := os.Mkdir(filepath.Join(cmdFolder, name), 0700); err != nil {
      return err
    }

    {
      //创键命令文件
      fileName := name + ".go"
      file := filepath.Join(cmdFolder, name, fileName)
      f, err := os.Create(file)
      if err != nil {
        return errors.Cause(err)
      }
      //创键title这个模板的方法
      funcs := template.FuncMap{"title": cases.Title(language.English).String}
      //使用contractTmp模板来初始化template,并且让这个模板支持title方法,即支持{{.|title}}
      t := template.Must(template.New("contract").Funcs(funcs).Parse(cmdTmp))
      if err := t.Execute(f, name); err != nil {
        return err
      }
    }

    fmt.Println("创建命令成功, 文件夹地址:", filepath.Join(cmdFolder, name))
    fmt.Println("请不要忘记挂载新创建的命令")

    return nil
  },
}

var cmdTmp string = `package {{.}}

import (
"github.com/SKYBroGardenLush/skycraper/framework/cobra"
)

var {{.|title}}Command = &cobra.Command{
  Use:     "{{.}}",
  Short:   "{{.}}的简要说明",
  Long:    "{{.}}的长说明",
  Example: "{{.}}命令的列子",
  RunE: func(cmd *cobra.Command, args []string) error {
    //书写你的逻辑
    return nil
  },
}
`
