package cobra

import (
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/robfig/cron/v3"
	"log"
)

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParantNull() {
	c.parent = nil
}

// AddCronCommand 是用来创键一个Cron任务的
func (c *Command) AddCronCommand(spec string, cmd *Command) {
	root := c.Root()
	if root.Cron == nil {
		//初始化cron
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional)))
		root.CronSpecs = []CronSpec{}

	}
	// 增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})
	//制作一个rootCommand
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParantNull()
	cronCmd.SetContainer(root.GetContainer())
	//增加调用函数
	root.Cron.AddFunc(spec, func() {
		//如果后续command出现panic,这里要捕获
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		fmt.Println("jjjj")
		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}
