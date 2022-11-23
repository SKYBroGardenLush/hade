package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/SKYBroGardenLush/skycraper/framework/cobra"
	"github.com/SKYBroGardenLush/skycraper/framework/survey"
	"github.com/SKYBroGardenLush/skycraper/framework/utils"
	"github.com/google/go-github/v48/github"
	"github.com/spf13/cast"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var newCommand = &cobra.Command{
	Use:   "new",
	Short: "脚手架命令",
	RunE: func(c *cobra.Command, args []string) error {

		var name string
		var folder string
		var mod string
		var version string
		var currentPath = utils.GetExecDirectory()
		var release *github.RepositoryRelease

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
			client := github.NewClient(nil)
			prompt := &survey.Input{
				Message: "请输入版本名称(参考 https://github.com/SKYBroGardenLush/skycraper/release,默认为最新版本):",
			}
			err := survey.AskOne(prompt, &version)
			if err != nil {
				return err
			}
			if version != "" {
				//确认版本号是否正确
				release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "SKYBroGardenLush", "skyscraper", version)
				if err != nil || release == nil {
					fmt.Println("所需版本不存在，创建失败，请参考 https://github.com/SKYBroGardenLush/skycraper/release")
				}
			} else if version == "" {
				release, _, err := client.Repositories.GetLatestRelease(context.Background(), "SKYBroGardenLush", "skyscraper")
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
				version = release.GetTagName()
			}
		}

		//==========下载和解压文件到相关文件夹===================
		templateFolder := filepath.Join(currentPath, "template-skyscraper-"+version+"-"+cast.ToString(time.Now().Unix()))
		os.Mkdir(templateFolder, os.ModePerm)
		fmt.Println("创建临时目录", templateFolder)

		//拷贝template项目
		url := release.GetZipballURL()
		err := utils.DownloadFile(filepath.Join(templateFolder, "template.zip"), url)
		if err != nil {
			return err
		}
		fmt.Println("下载zip包到template.zip")
		_, err = utils.Unzip(filepath.Join(templateFolder, "template.zip"), templateFolder)
		if err != nil {
			return err
		}

		//获取folder下的SKYBroGardenLush-skyscraper-xxx 相关目录
		fInfos, err := ioutil.ReadDir(templateFolder)
		if err != nil {
			return err
		}
		for _, fInfo := range fInfos {
			//找到解开后的文件夹
			if fInfo.IsDir() && strings.Contains(fInfo.Name(), "SKYBroGardenLush-skyscraper-") {
				if err := os.Rename(filepath.Join(templateFolder, fInfo.Name()), folder); err != nil {
					if err != nil {
						return err
					}
				}
			}
		}
		fmt.Println("解压zip包")

		if err := os.RemoveAll(templateFolder); err != nil {
			return err
		}
		fmt.Println("删除临时文件夹")
		//==========完成 下载和解压文件到相关文件夹===================

		//========== 修改文件相应信息和删除一些不必要的文件 =============================
		os.RemoveAll(path.Join(folder, "git"))
		fmt.Println("删除.git目录")

		//删除framework 目录
		os.RemoveAll(path.Join(folder, "framework"))
		fmt.Println("删除framework目录")

		filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			//修改go.mod中的模块名称，修改go.mod中的require信息
			//增加require github.com/SKYBroGardenLush/skyscraper
			if path == filepath.Join(folder, "go.mod") {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("module github.com/SKYBroGardenLush/skyscraper"), []byte("module"+mod))
				c = bytes.ReplaceAll(c, []byte("require("), []byte("require (\n\tgithub.com/SKYBroGardenLush/skyscraper "+version))
				err := ioutil.WriteFile(path, c, 0664)
				if err != nil {
					return err
				}
				return nil
			}

			isContain := bytes.Contains(c, []byte("github.com/SKYBroGardenLush/skyscraper/app"))
			if isContain {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("github.com/SKYBroGardenLush/skyscraper/app"), []byte(mod+"app"))
				err = ioutil.WriteFile(path, c, 0664)
				if err != nil {
					return err
				}
			}
			return nil

		})
		fmt.Println("创建成功！！")
		fmt.Println("目录：", folder)
		fmt.Println("=====================================")
		return nil
	},
}
