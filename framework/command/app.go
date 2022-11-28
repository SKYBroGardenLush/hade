package command

import (
	"context"
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/cobra"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
	"github.com/SKYBroGardenLush/skyscraper/framework/utils"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		//打印帮助文档
		cmd.Help()
		return nil
	},
}

//app 启动地址
var appAddress = ""
var host = ""
var port = ""
var appDaemon = false

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		//从Command中获取服务容器
		container := cmd.Root().GetContainer()
		//检查app.pid和app.log记录文件是否存在
		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		if !utils.Exists(pidFolder) {
			err := os.MkdirAll(pidFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !utils.Exists(logFolder) {
			err := os.MkdirAll(logFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		serverLogFile := filepath.Join(pidFolder, "app.log")
		currentFolder := utils.GetExecDirectory()

		//从服务容器中获取kernel的服务实例
		kernelServie := container.MustMake(contract.KernelKey).(contract.Kernel)
		//从kernel服务实例中获取引擎
		core := kernelServie.HttpEngine()
		//获取监听地址

		if host == "127.0.0.1" && port == "8080" { //为默认地址，则读取配置文件
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8080"
				}
			}
		} else {
			appAddress = host + ":" + port
		}

		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}
		//后台启动
		if appDaemon {
			//创建一个context
			cntxt := &daemon.Context{
				//设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				//设置日志文件
				LogFileName: serverLogFile,
				LogFilePerm: 0664,
				//设置工作路径
				WorkDir: currentFolder,
				//设置所有文件的mask默认为750
				Umask: 027,
				//子进程参数,子进程命令为 ./hade app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			//启动子进程，d为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				//父进程打印启动消息，不做操作
				fmt.Println("app 启动成功，pid:", d.Pid)
				fmt.Println("日志文件", serverLogFile)
				return nil
			}
			defer cntxt.Release()
			//子进程执行真正的app启动操作
			fmt.Println("deamon started")
			gspt.SetProcTitle("hade app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		}

		//非deamon 模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("hade app")
		err = startAppServe(server, container)
		if err != nil {
			return err
		}
		return nil
	},
}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid,暂不支持windows操作系统",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.Root().GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		//获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}
		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))

			if err != nil {
				return err
			}
			//检查pid是否存在
			if utils.CheckProcessExist(pid) {
				fmt.Println("app 服务已经启动，pid:", pid)
			}
		}
		fmt.Println("app 没有服务存在")
		return nil
	},
}

//var appStopCommand = &cobra.Command{
//	Use:   "stop",
//	Short: "停止一个已经启动的app服务",
//	RunE: func(c *cobra.Command, args []string) error {
//		container := c.GetContainer()
//		appService := container.MustMake(contract.AppKey).(contract.App)
//
//		//获取pid
//		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
//		content, err := ioutil.ReadFile(serverPidFile)
//		if err != nil {
//			return err
//		}
//		if content != nil && len(content) > 0 {
//			pid, err := strconv.Atoi(string(content))
//			if err != nil {
//				return err
//			}
//			//发送SIGTERM命令
//			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
//				return err
//			}
//			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
//				return err
//			}
//
//		}
//		return nil
//	},
//}

//// 重新启动一个app服务
//var appRestartCommand = &cobra.Command{
//	Use:   "restart",
//	Short: "重新启动一个app服务",
//	RunE: func(c *cobra.Command, args []string) error {
//		container := c.GetContainer()
//		appService := container.MustMake(contract.AppKey).(contract.App)
//
//		// GetPid
//		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
//
//		content, err := ioutil.ReadFile(serverPidFile)
//		if err != nil {
//			return err
//		}
//
//		if content != nil && len(content) != 0 {
//			pid, err := strconv.Atoi(string(content))
//			if err != nil {
//				return err
//			}
//			if utils.CheckProcessExist(pid) {
//				// 杀死进程
//				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
//					return err
//				}
//				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
//					return err
//				}
//
//				// 获取closeWait
//				closeWait := 5
//				configService := container.MustMake(contract.ConfigKey).(contract.Config)
//				if configService.IsExist("app.close_wait") {
//					closeWait = configService.GetInt("app.close_wait")
//				}
//
//				// 确认进程已经关闭,每秒检测一次， 最多检测closeWait * 2秒
//				for i := 0; i < closeWait*2; i++ {
//					if utils.CheckProcessExist(pid) == false {
//						break
//					}
//					time.Sleep(1 * time.Second)
//				}
//
//				// 如果进程等待了2*closeWait之后还没结束，返回错误，不进行后续的操作
//				if utils.CheckProcessExist(pid) == true {
//					fmt.Println("结束进程失败:"+strconv.Itoa(pid), "请查看原因")
//					return errors.New("结束进程失败")
//				}
//
//				fmt.Println("结束进程成功:" + strconv.Itoa(pid))
//			}
//		}
//
//		appDaemon = true
//		// 直接daemon方式启动apps
//		return appStartCommand.RunE(c, args)
//	},
//}

func startAppServe(server *http.Server, c framework.Container) error {
	go func() {
		fmt.Println("服务地址:", server.Addr)
		server.ListenAndServe()
	}()

	//当前Goroutine 等待信号量
	quit := make(chan os.Signal)
	// 监控信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前 Goroutine 等待信号
	<-quit

	//调用server.shutdown 等待所有连接处理完毕
	closeWaitTime := 5
	configServer := c.MustMake(contract.ConfigKey).(contract.Config)
	if configServer.IsExist("app.close_wait") {
		closeWaitTime = configServer.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWaitTime)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Println("server shutdown: ", err.Error())
	}
	return nil
}

//初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "设置监听主机地址,默认127.0.0.1")
	appStartCommand.Flags().StringVarP(&port, "port", "P", "8080", "设置监听主机地址，默认8080")
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "是否后台启动,默认前台启动,暂不支持Non-POSIX OS(如windows操作系统)")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	//appCommand.AddCommand(appStopCommand)
	//appCommand.AddCommand(appRestartCommand)
	return appCommand
}
