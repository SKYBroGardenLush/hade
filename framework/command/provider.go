package command

import (
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/cobra"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
	"github.com/SKYBroGardenLush/skyscraper/framework/survey"
	"github.com/SKYBroGardenLush/skyscraper/framework/utils"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"text/template"
)

func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerCreateCommand)
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供者相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "服务提供者相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		providers := c.GetContainer().(*framework.HadeContainer).NameList() //注意这个*号
		//打印字符凭证
		for _, provider := range providers {
			fmt.Println(provider)
		}
		return nil
	},
}

var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "服务提供者相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("创键一个服务")
		var name string
		var folder string
		//服务凭证
		{
			prompt := &survey.Input{
				Message: "请输入服务名称(服务凭证)",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		//服务名称
		{
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称(默认:同服务名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		//检查服务是否存在
		providers := container.(*framework.HadeContainer).NameList()
		providerColl := utils.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务已经存在")
			return nil
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder := app.ProviderFolder()
		subFolders, err := utils.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := utils.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		//开始创键文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		//创键title这个模板的方法
		funcs := template.FuncMap{"title": cases.Title(language.English).String}
		{
			//创键contract.go
			file := filepath.Join(pFolder, folder, "contact.go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			//使用contractTmp模板来初始化template,并且让这个模板支持title方法,即支持{{.|title}}
			t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}

		{
			// 创建provider.go
			file := filepath.Join(pFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("provider").Funcs(funcs).Parse(providerTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		{
			//  创建service.go
			file := filepath.Join(pFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("service").Funcs(funcs).Parse(serviceTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		fmt.Println("创建服务成功, 文件夹地址:", filepath.Join(pFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")

		return nil
	},
}

var contractTmp string = `package {{.}}

const {{.|title}}Key = "{{.}}"

type Service interface{
  //请在这里定义你的方法
  Foo() string
}
`

var providerTmp string = `package {{.}}
import (
	"github.com/SKYBroGardenLush/skyscraper/framework"
)
type {{.|title}}Provider struct {
	framework.ServiceProvider
	c framework.Container
}
func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}
func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}
func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}
func (sp *{{.|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	return nil
}
`

var serviceTmp string = `package {{.}}
import "github.com/SKYBroGardenLush/skyscraper/framework"
type {{.|title}}Service struct {
	container framework.Container
}
func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container}, nil
}
func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
