package contract

import (
	"context"
	"io"
	"time"
)

const LogKey = "hade:log"

type LogLevel uint32

const (
	// UnknownLevel 表示未知的日志级别
	UnknownLevel LogLevel = iota
	// PanicLevel 表示导致整个程序出现崩溃的日志信息
	PanicLevel
	// FatalLevel 表示会导致当前这个请求出现提前终止的错误信息
	FatalLevel
	// ErrorLevel 表示出现错误但不一定影响后续请求逻辑的错误信息
	ErrorLevel
	// WarnLevel 表示出现错误但不一定影响后续请求逻辑的警报信息
	WarnLevel
	// InfoLevel 表示正常的日志信息输出
	InfoLevel
	// DebugLevel 表示调试状态下打印出来的日志信息
	DebugLevel
	// TraceLevel 表示最详细的信息，可能包含堆栈等信息
	TraceLevel
)

type Log interface {
	// Panic 表示导致整个程序出现崩溃的日志信息
	Panic(ctx context.Context, msg string, fields map[string]interface{})
	// Fatal 表示会导致当前这个请求出现提前终止的错误信息
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
	// Error 表示出现错误但不一定影响后续请求逻辑的错误信息
	Error(ctx context.Context, msg string, fields map[string]interface{})
	// Warn 表示出现错误但不一定影响后续请求逻辑的警报信息
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	// Info 表示正常的日志信息输出
	Info(ctx context.Context, msg string, fields map[string]interface{})
	// Debug 表示调试状态下打印出来的日志信息
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	// Trace 表示最详细的信息，可能包含堆栈等信息
	Trace(ctx context.Context, msg string, fields map[string]interface{})
	// SetLevel 设置日志级别
	SetLevel(level LogLevel)
	// SetCtxFielder 设置从context中获取上下文字段field方法
	SetCtxFielder(handler CtxFielder)
	// SetFormatter 设置输出格式
	SetFormatter(formatter Formatter)
	// SetOutput 设置输出管道
	SetOutput(out io.Writer)
}

// CtxFielder 定义从context中获取信息的方法
type CtxFielder func(ctx context.Context) map[string]interface{}

type Formatter func(level LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error)
