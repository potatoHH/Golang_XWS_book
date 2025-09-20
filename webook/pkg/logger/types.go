package logger

// TODO 风格1
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// TODO 风格2
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	//args 会加入加入进去任何LoggerV1中的任何打印出来的日志中
	With(args ...Field) LoggerV1
}
type Field struct {
	Key   string
	Value any
}

// TODO 风格3  要求args 为偶数
type LoggerV2 interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// TODO 对风格2 进行封装
