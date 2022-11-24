package command

import (
	"errors"
	"fmt"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/cobra"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
	"github.com/SKYBroGardenLush/skyscraper/framework/utils"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Proxy struct {
	devConfig   *devConfig //配置文件
	backendPid  int        //当前backend服务的pid
	frontendPid int        //当前frontend服务的pid
}

func NewProxy(c framework.Container) *Proxy {
	return &Proxy{}
}

//重新启动一个proxy 网关
func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 && p.backendPid == 0 {
		fmt.Println("前端和后端服务都不存在")
		return nil
	}

	//后端服务存在
	if p.frontendPid == 0 && p.backendPid != 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	}

	//前端服务存在
	if p.frontendPid != 0 && p.backendPid == 0 {
		return httputil.NewSingleHostReverseProxy(frontend)
	}

	//两个进程都有
	//先创建一个后端的directory
	director := func(req *http.Request) {
		if req.URL.Path == "/" || req.URL.Path == "/app.js" {
			req.URL.Scheme = frontend.Scheme
			req.URL.Host = frontend.Host
		} else {
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
		}
	}

	//定义一个NotFoundErr
	NotFoundErr := errors.New("response is 404,need to redirect")
	return &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(response *http.Response) error {
			//如果后端服务返回了404，我们返回NotFoundErr 会进入到errorHandler中
			if response.StatusCode == 404 {
				return NotFoundErr
			}
			return NotFoundErr
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			//判断Error是否为NotFoundError,是的话则进行前端服务的转发，重新修改writer
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(frontend).ServeHTTP(writer, request)
			}
		},
	}
}

func (p *Proxy) startProxy(startFrontend, startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error

	//启动后端
	if startBackend {
		if err = p.restartBackend(); err != nil {
			return err
		}
	}
	//启动前端
	if startFrontend {
		if err = p.restartFrontend(); err != nil {
			return err
		}
	}

	if frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Frontend.Port)); err != nil {
		return err
	}

	if frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port)); err != nil {
		return err
	}

	//设置反向代理
	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	proxyServer := &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}

	fmt.Println("代理服务启动:", "http://"+proxyServer.Addr)
	err = proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//启动前端服务
func (p *Proxy) restartFrontend() error {
	//启动前端调试模式
	//先杀死就进程
	if p.frontendPid != 0 {
		//syscall.Kill(p.frontendPid, syscall.SIGKILL)
		p.frontendPid = 0
	}
	//否则开启 npm run serve
	port := p.devConfig.Frontend.Port
	path, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "run", "dev")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s%s", "PORT=", port))
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	//因为npm run serve 是控制台挂起模式，所以这里使用go routine 启动
	err = cmd.Start()
	fmt.Println("启动前端服务:", "http://127.0.0.1"+port)
	if err != nil {
		fmt.Println(err)
	}
	p.frontendPid = cmd.Process.Pid
	fmt.Println("前端服务pid:", p.frontendPid)
	return nil
}

func (p *Proxy) restartBackend() error {
	//杀死之前的进程
	if p.backendPid != 0 {
		//syscall.Kill(p.backendPid, syscall.SIGKILL)
		p.backendPid = 0
	}
	//设置随机端口，真实后端端口
	port := p.devConfig.Backend.Port
	hadeAddress := fmt.Sprintf(":" + port)
	//使用命令行启动后端进程
	cmd := exec.Command("./hade", "app", "start", "--address="+hadeAddress)
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务:", "http://127.0.0.1:"+port)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务pid:", p.backendPid)
	return nil
}
func (p *Proxy) rebuildBackend() error {
	//重新编译hade后端
	cmdBuild := exec.Command("./hade", "build", "backend")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Start(); err == nil {
		err = cmdBuild.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Proxy) monitorBackend() error {
	//监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	//开启监听目标文件夹
	appFolder := p.devConfig.Backend.MonitorFolder
	fmt.Println("监控文件夹:", appFolder)
	//监听所有子目录，需要使用filepath.walk
	filepath.Walk(appFolder, func(path string, info fs.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		//如果是隐藏目录则不用监听
		if utils.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})

	//开启计时时间机制
	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	//先停止计时器
	t.Stop()
	for {
		select {
		case <-t.C:
			//记时时间到了，代表之前有文件更新事件重置过时计时器
			//即有文件更新
			fmt.Println("...检测到文件更新，重新服务开始...")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译后端失败", err.Error())
			} else {
				if err := p.restartFrontend(); err != nil {
					fmt.Println("重新启动后端失败：", err.Error())
				}
			}
			fmt.Println("...检测到文件更新,重启服务结束...")
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			//如果有文件监听错误，则停止计时器
			fmt.Println("监听文件夹错误", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}

	}

}

type devConfig struct {
	Port string //调试模式最终监听的端口，默认为8070

	Backend struct { //后端调试模式配置
		RefreshTime   int    //调试模式后端更新时间，如果文件变更，等待3s才进行一次更新
		Port          string //后端监听端口，默认为8072
		MonitorFolder string //监听文件夹，默认为AppFolder
	}

	Frontend struct { //前端调试模式配置
		Port string //前端启动端口，默认8071
	}
}

func initDevConfig(c framework.Container) *devConfig {
	//设置默认值
	devConfig := &devConfig{
		Port: "8087",
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{RefreshTime: 1, Port: "8072", MonitorFolder: ""},

		Frontend: struct {
			Port string
		}{
			Port: "8071",
		},
	}
	//容器中获取配置服务
	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	// 每个配置项进行检查
	if configer.IsExist("app.dev.port") {
		devConfig.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExist("app.dev.backend.refresh_time") {
		devConfig.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	if configer.IsExist("app.dev.backend.port") {
		devConfig.Port = configer.GetString("app.dev.backend.port")
	}

	//monitorFolder 默认使用目录服务的AppFolder()
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := c.MustMake(contract.AppKey).(contract.App)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	}
	if configer.IsExist("app..dev.frontend.port") {
		devConfig.Frontend.Port = configer.GetString("app.dev.frontend.port")
	}

	return devConfig
}

func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devAllCommand)
	devCommand.AddCommand(devFrontendCommand)
	devCommand.AddCommand(devBackendCommand)
	return devCommand
}

var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var devAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时启动前后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		//启动前端和后端服务
		proxy := NewProxy(c.GetContainer())
		//监听后端文件
		go proxy.monitorBackend()
		//启动只有后端的proxy
		if err := proxy.startProxy(true, true); err != nil {
			return err
		}
		return nil
	},
}

var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "启动后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {

		//启动后端服务
		proxy := NewProxy(c.GetContainer())
		//监听后端文件
		go proxy.monitorBackend()
		//启动只有后端的proxy
		if err := proxy.startProxy(false, true); err != nil {
			return err
		}
		return nil
	},
}

var devFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "启动前端调试模式",
	RunE: func(c *cobra.Command, args []string) error {

		//启动前端服务
		proxy := NewProxy(c.GetContainer())
		return proxy.startProxy(true, false)
	},
}
