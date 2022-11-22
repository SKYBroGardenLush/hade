package contract

const (
	// EnvProduction 代表生产环境
	EnvProduction = "production"
	// EnvTesting 代表测试环境
	EnvTesting = "testing"
	// EnvDevelopment 代表开发环境
	EnvDevelopment = "development"
	// EnvKey 环境变量服务凭证
	EnvKey = "hade:env"
)

// Env 定义环境变量服务
type Env interface {
	// AppEnv 获取当前环境变量
	AppEnv() string
	// IsExit 判断一个环境变量是否被设置
	IsExit(string) bool
	// Get 获取某个环境变量，如果没有则返回""
	Get(string) string
	// All 获取所有环境变量，.env和运行环境变量融合的结果
	All() map[string]string
}
