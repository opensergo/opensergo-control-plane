package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a thin wrapper for zap.Logger that adds Ctx method
type ZapLogger struct {
	*zap.Logger
}

func New(o *Options) *ZapLogger {
	if o == nil {
		o = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       false,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "timestamp",
			NameKey:    "logger",
			CallerKey:  "caller",
			//FunctionKey:    "func", // if need
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			//EncodeName:       zapcore.FullNameEncoder,
			//ConsoleSeparator: "",
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
		//InitialFields:    nil,
	}
	logger, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(err)
	}
	//logger = logger.Named(o.Name)

	//zap.RedirectStdLog(logger.Named(o.Name))
	//zap.ReplaceGlobals(logger)

	return &ZapLogger{
		Logger: logger,
	}
}

// Flush calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Flush() error {
	return defaultLog.Sync()
}

func WithName(s string) *ZapLogger {
	return defaultLog.WithName(s)
}

// WithName adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (zl *ZapLogger) WithName(name string) *ZapLogger {
	return &ZapLogger{
		zl.Named(name),
	}
}

func WithValues(keysAndValues ...interface{}) *ZapLogger {
	return defaultLog.WithValues(keysAndValues...)
}

func (zl *ZapLogger) WithValues(keysAndValues ...interface{}) *ZapLogger {
	return &ZapLogger{
		zl.With(zl.handleFields(keysAndValues)...),
	}
}

func (zl *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	zl.Logger.Info(msg, zl.handleFields(keysAndValues)...)
}

func (zl *ZapLogger) Infof(format string, args ...interface{}) {
	zl.Logger.Sugar().Infof(format, args...)
}

func (zl *ZapLogger) InfoC(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// todo
	//keysAndValues = zl.FromBizContext(ctx, keysAndValues...)
	zl.Logger.Info(msg, zl.handleFields(keysAndValues)...)
}

func (zl *ZapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "log", zl)
}

// handleFields converts a bunch of arbitrary key-value pairs into Zap fields.  It takes
// additional pre-converted Zap fields, for use with automatically attached fields, like
// `error`.
func (zl *ZapLogger) handleFields(args []interface{}) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// Slightly slower fast path when we need to inject "v".
		return []zap.Field{}
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args); {
		// Check just in case for strongly-typed Zap fields,
		// which might be illegal (since it breaks
		// implementation agnosticism). If disabled, we can
		// give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			zl.WithOptions(zap.AddCallerSkip(1)).DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			zl.WithOptions(zap.AddCallerSkip(1)).DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			zl.WithOptions(zap.AddCallerSkip(1)).DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields)
}
