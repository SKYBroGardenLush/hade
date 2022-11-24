package main

import (
	"github.com/SKYBroGardenLush/skyscraper/app/console"
	"github.com/SKYBroGardenLush/skyscraper/app/http"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/provider/app"
	"github.com/SKYBroGardenLush/skyscraper/framework/provider/config"
	"github.com/SKYBroGardenLush/skyscraper/framework/provider/distributed"
	"github.com/SKYBroGardenLush/skyscraper/framework/provider/env"
	"github.com/SKYBroGardenLush/skyscraper/framework/provider/kernel"
)

func main() {
	//core := framework.NewCore()
	//registerRouter(core)

	//core := gin.New()
	//
	////绑定具体的服务
	//core.Bind(&demo.DemoServiceProvider{})
	//registerRouterGin(core)
	//server := &http.Server{
	//
	//	Handler: core,
	//	Addr:    ":8081",
	//}
	//go func() {
	//	server.ListenAndServe()
	//}()
	//
	////当前Goroutine 等待信号量
	//quit := make(chan os.Signal)
	//// 监控信号
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//// 这里会阻塞当前 Goroutine 等待信号
	//<-quit
	//
	////调用server.shutdown 等待所有连接处理完毕
	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//if err := server.Shutdown(timeoutCtx); err != nil {
	//	fmt.Println("server shutdown: ", err.Error())
	//}

	container := framework.NewHeadContainer()

	container.Bind(&app.HadeAppProvider{})

	// 后续初始化需要绑定的服务提供者...
	container.Bind(&distributed.LocalDistributedProvider{})

	container.Bind(&env.HadeEnvProvider{})

	container.Bind(&config.HadeConfigProvider{})

	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)

}
