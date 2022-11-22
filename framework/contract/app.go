package contract

const AppKey = "hade:app"

type App interface {
	// AppID 表示当前app唯一id可用于分布式锁等
	AppID() string
	// Version 定义当前版本
	Version() string
	// BaseFolder 定义项目基础地址
	BaseFolder() string
	//AppFolder 定义项目app地址
	AppFolder() string
	// ConfigFolder 定义配置文件路径
	ConfigFolder() string
	// LogFolder 定义日志文件路径
	LogFolder() string
	// ProviderFolder 定义业务自己的服务提供者地址
	ProviderFolder() string
	// MiddlewareFolder 定义业务自己定义的中间件
	MiddlewareFolder() string
	// CommandFolder 定义业务的命令
	CommandFolder() string
	// RuntimeFolder 定义业务运行的中间态信息
	RuntimeFolder() string
	// TestFolder 存放测试所需要的信息
	TestFolder() string
	LoadAppConfig(map[string]string)
}
