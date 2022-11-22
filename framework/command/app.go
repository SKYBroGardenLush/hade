package command

import (
  "context"
  "fmt"
  "github.com/SKYBroGardenLush/skycraper/framework/cobra"
  "github.com/SKYBroGardenLush/skycraper/framework/contract"
  "net/http"
  "os"
  "os/signal"
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

var appStartCommand = &cobra.Command{
  Use:   "start",
  Short: "启动一个web服务",
  RunE: func(cmd *cobra.Command, args []string) error {
    //从Command中获取服务容器
    container := cmd.Root().GetContainer()
    //从服务容器中获取kernel的服务实例
    kernelServie := container.MustMake(contract.KernelKey).(contract.Kernel)
    //从kernel服务实例中获取引擎
    core := kernelServie.HttpEngine()
    server := &http.Server{
      Handler: core,
      Addr:    appAddress,
    }
    go func() {
      server.ListenAndServe()
    }()

    //当前Goroutine 等待信号量
    quit := make(chan os.Signal)
    // 监控信号
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    // 这里会阻塞当前 Goroutine 等待信号
    <-quit

    //调用server.shutdown 等待所有连接处理完毕
    timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(timeoutCtx); err != nil {
      fmt.Println("server shutdown: ", err.Error())
    }
    return nil
  },
}

//初始化app命令和其子命令
func initAppCommand() *cobra.Command {
  appStartCommand.Flags().StringVar(&appAddress, "address", ":8888", "设置app启动地址默认为:8888")
  appCommand.AddCommand(appStartCommand)
	return appCommand
}
