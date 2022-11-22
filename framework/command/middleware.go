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
  "io/ioutil"
  "os"
  "path/filepath"
  "text/template"
)

func initMiddlewareCommand() *cobra.Command {
  middlewareCommand.AddCommand(middlewareListCommand)
  middlewareCommand.AddCommand(middlewareCreateCommand)
  return middlewareCommand
}

var middlewareCommand = &cobra.Command{
  Use:   "middleware",
  Short: "中间件相关命令",
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.Help()
    return nil
  },
}

var middlewareListCommand = &cobra.Command{
  Use:   "list",
  Short: "显示所有安装的中间件",
  RunE: func(cmd *cobra.Command, args []string) error {
    contrainer := cmd.Root().GetContainer()
    appService := contrainer.MustMake(contract.AppKey).(contract.App)
    middlewarePath := filepath.Join(appService.BaseFolder(), "app", "http", "middleware")
    // 读取文件夹
    files, err := ioutil.ReadDir(middlewarePath)
    if err != nil {
      return err
    }
    for _, f := range files {
      if f.IsDir() {
        fmt.Println(f.Name())
      }
    }
    return nil
  },
}

var middlewareCreateCommand = &cobra.Command{
  Use:     "new",
  Aliases: []string{"create", "init"},
  Short:   "创键中间件",
  RunE: func(c *cobra.Command, args []string) error {
    container := c.GetContainer()
    fmt.Println("创键一个中间件")
    var name string
    var folder string
    //服务凭证
    {
      prompt := &survey.Input{
        Message: "请输入中间件文件名称:",
      }
      err := survey.AskOne(prompt, &name)
      if err != nil {
        return err
      }
    }
    //服务名称
    {
      prompt := &survey.Input{
        Message: "请输入中间件所在目录名称(默认:同服务名称):",
      }
      err := survey.AskOne(prompt, &folder)
      if err != nil {
        return err
      }
    }

    //检查服务是否存在

    if folder == "" {
      folder = name
    }

    app := container.MustMake(contract.AppKey).(contract.App)

    middleFolder := app.MiddlewareFolder()
    subFolders, err := utils.SubDir(middleFolder)
    if err != nil {
      return err
    }
    subColl := utils.NewStrCollection(subFolders)
    if subColl.Contains(folder) {
      fmt.Println("目录名称已经存在")
    } else {
      //目录不存在创键新的目录
      if err := os.Mkdir(filepath.Join(middleFolder, folder), 0700); err != nil {
        return err
      }
      fmt.Println("成功创键新目录:", filepath.Join(middleFolder, folder))
    }

    //创键新文件
    subFiles, err := utils.SubFile(middleFolder)
    if err != nil {
      return err
    }
    subColl = utils.NewStrCollection(subFiles)
    fileName := name + ".go"
    if subColl.Contains(fileName) {
      fmt.Println("文件已经存在")
    } else {
      //创键这个文件
      file := filepath.Join(middleFolder, folder, fileName)
      f, err := os.Create(file)
      if err != nil {
        return errors.Cause(err)
      }
      //创键title这个模板的方法
      funcs := template.FuncMap{"title": cases.Title(language.English).String}
      //使用contractTmp模板来初始化template,并且让这个模板支持title方法,即支持{{.|title}}
      t := template.Must(template.New("contract").Funcs(funcs).Parse(middlewareTmp))
      if err := t.Execute(f, folder); err != nil {
        return err
      }

    }

    fmt.Println("创建中间件文件成功, 文件地址:", filepath.Join(middleFolder, folder, fileName))

    return nil
  },
}

//// 从gin-contrib中迁移中间件
//var middlewareMigrateCommand = &cobra.Command{
//	Use:   "migrate",
//	Short: "迁移gin-contrib中间件, 迁移地址：https://github.com/gin-contrib/[middleware].git",
//	RunE: func(c *cobra.Command, args []string) error {
//		container := c.GetContainer()
//		fmt.Println("迁移一个Gin中间件")
//		var repo string
//		{
//			prompt := &survey.Input{
//				Message: "请输入中间件名称：",
//			}
//			err := survey.AskOne(prompt, &repo)
//			if err != nil {
//				return err
//			}
//		}
//		// step2 : 下载git到一个目录中
//		appService := container.MustMake(contract.AppKey).(contract.App)
//
//		middlewarePath := appService.MiddlewareFolder()
//		url := "https://github.com/gin-contrib/" + repo + ".git"
//		fmt.Println("下载中间件 gin-contrib:")
//		fmt.Println(url)
//		_, err := git.PlainClone(path.Join(middlewarePath, repo), false, &git.CloneOptions{
//			URL:      url,
//			Progress: os.Stdout,
//		})
//		if err != nil {
//			return err
//		}
//
//		// step3:删除不必要的文件 go.mod, go.sum, .git
//		repoFolder := path.Join(middlewarePath, repo)
//		fmt.Println("remove " + path.Join(repoFolder, "go.mod"))
//		os.Remove(path.Join(repoFolder, "go.mod"))
//		fmt.Println("remove " + path.Join(repoFolder, "go.sum"))
//		os.Remove(path.Join(repoFolder, "go.sum"))
//		fmt.Println("remove " + path.Join(repoFolder, ".git"))
//		os.RemoveAll(path.Join(repoFolder, ".git"))
//
//		// step4 : 替换关键词
//		filepath.Walk(repoFolder, func(path string, info os.FileInfo, err error) error {
//			if info.IsDir() {
//				return nil
//			}
//
//			if filepath.Ext(path) != ".go" {
//				return nil
//			}
//
//			c, err := ioutil.ReadFile(path)
//			if err != nil {
//				return err
//			}
//			isContain := bytes.Contains(c, []byte("github.com/gin-gonic/gin"))
//			if isContain {
//				fmt.Println("更新文件:" + path)
//				c = bytes.ReplaceAll(c, []byte("github.com/gin-gonic/gin"), []byte("github.com/SKYBroGardenLush/skycraper/framework/gin"))
//				err = ioutil.WriteFile(path, c, 0644)
//				if err != nil {
//					return err
//				}
//			}
//
//			return nil
//		})
//		return nil
//	},
//}

var middlewareTmp string = `package {{.}}
import "github.com/SKYBroGardenLush/skycraper/framework/gin"
// {{.|title}}Middleware 代表中间件函数
func {{.|title}}Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
    //您的处理代码
		context.Next()
	}
}
`
