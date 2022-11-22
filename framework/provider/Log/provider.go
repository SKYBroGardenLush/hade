package Log

import (
	"github.com/SKYBroGardenLush/skycraper/framework"
	"github.com/SKYBroGardenLush/skycraper/framework/contract"
	"github.com/SKYBroGardenLush/skycraper/framework/provider/Log/formatter"
	"github.com/SKYBroGardenLush/skycraper/framework/provider/Log/services"

	"io"
	"strings"
)

type HadeLogServiceProvider struct {
	//选择服务
	Driver string
	//日志级别
	Level contract.LogLevel
	//日志格式
	Formatter contract.Formatter
	//日志上下文信息获取函数
	CtxFielder contract.CtxFielder
	//日志输出信息
	Output io.Writer
}

func (l *HadeLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			//默认使用console
			return services.NewHadeConsoleLog
		}
		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}
	//根据driver的配置项确定
	switch l.Driver {
	case "single":
		return services.NewHadeSingleLog
	case "custom":
		return services.NewHadeCustomLog
	case "console":
		return services.NewHadeConsoleLog
	case "rotate":
		return services.NewHadeRotateLog
	default:
		return services.NewHadeConsoleLog

	}
}

func (l *HadeLogServiceProvider) Params(c framework.Container) []interface{} {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	//设置formatter
	if l.Formatter == nil {
		l.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				l.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				l.Formatter = formatter.TextFormatter
			}
		}

	}
	if l.Level == contract.UnknownLevel {
		l.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			l.Level = logLevel(configService.GetString("log.level"))
		}

	}
	//定义五个参数
	return []interface{}{c, l.Level, l.CtxFielder, l.Formatter, l.Output}

}

func (l *HadeLogServiceProvider) Name() string {
	return contract.LogKey
}

func (l *HadeLogServiceProvider) IsDefer() bool {
	return false
}

func logLevel(level string) contract.LogLevel {
	switch level {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	default:
		return contract.InfoLevel

	}

}
