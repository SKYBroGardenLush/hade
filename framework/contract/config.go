package contract

import "time"

const ConfigKey = "hade:config"

// Config 定义了配置文件服务，读取配置文件，支持点分割的路径读取
//建议使用yaml属性
type Config interface {
	IsExist(key string) bool
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetString(key string) string
	GetStringSlice(key string) []string
	GetIntSlice(key string) []int
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	// Load 加载配置到某个对象
	Load(key string, val interface{}) error
}
