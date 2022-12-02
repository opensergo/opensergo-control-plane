package log

import (
	"context"
)

func init() {
	defaultLog = New(nil)
}

var defaultLog = New(NewOptions())

func FromContext(ctx context.Context) *ZapLogger {
	if ctx != nil {
		logger := ctx.Value("log")
		if logger != nil {
			return logger.(*ZapLogger)
		}
	}
	return WithName("Unknown-Context")
}

func Info(msg string, keysAndValues ...interface{}) {
	defaultLog.Info(msg, keysAndValues...)
}

func Infof(format string, args ...interface{}) {
	defaultLog.Infof(format, args...)
}

func InfoC(ctx context.Context, msg string, keysAndValues ...interface{}) {
	defaultLog.InfoC(ctx, msg, keysAndValues...)
}
