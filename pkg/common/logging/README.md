# logging component

here provider a common Logger interface which is highly extensible and multi implementable

## Quick start
In this logging component, you can have only one global ConsoleLogger, and several FileLogger, By invoking function in [logger.go](./logger.go):
``` go
    // AppendLoggerSlice add the Logger into loggerSlice
    func AppendLoggerSlice(loggerAppend Logger) {
        loggerSlice = append(loggerSlice, loggerAppend)
    }
    
    // SetConsoleLogger set the consoleLogger to print int stdout
    func SetConsoleLogger(logger Logger) {
        consoleLogger = logger
    }
```

You can use the default implement of Logger directly, By new Logger instance implemented in [logger_default.go](./logger_default.go),  
But first, you should new instances from [logger_default.go](./logger_default.go), like following: 
``` go
    // new instances
    logging.NewDefaultConsoleLogger(logging.DebugLevel)
    logging.NewDefaultFileLogger(logging.DebugLevel)
    // invoke the function in Logger interface
    logging.Error(errors.New("errors.New"), "this is error log in printErrorStack()")
    logging.Info("this is info log in printErrorStack()")
    logging.Debug("this is debug log in printErrorStack()")
```

You can format you log, by using [logger_assembler.go](./logger_assembler.go),  
the assembler provides  AssembleMsgJsonFormat and AssembleMsgSeparateFormat can be chosen.

``` go
    // the unified entrance for AssembleMsg
    func AssembleMsg(logFormat LogFormat, callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{}) string
    
    // AssembleMsgJsonFormat Assemble log-msg as json-format
    //
    // debug/info/warn:
    // {"logLevel":"INFO","timestamp":"2006-01-02 15:04:05.000","caller":"opensergo_client.go:74","msg":"openSergoClient is starting..."}
    //
    // error:
    // {"logLevel":"ERROR","timestamp":"2006-01-02 15:04:05.000","caller":"opensergo_client.go:83","msg":"can not connect.","error":"rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp 33.1.33.1:10246: connect: connection refused\""}
    // [ ERROR CAUSES WITH STACK ][ rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 33.1.33.1:10246: connect: connection refused"
    // github.com/opensergo/opensergo-go/pkg/client.handleReceive.func1
    //
    //	.../opensergo-go-sdk/pkg/client/opensergo_client.go:82 ][ ERROR CAUSES WITH STACK END ]
    func AssembleMsgJsonFormat(callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{}) string
    
    // AssembleMsgSeparateFormat Assemble log-msg as separate-format
    //
    // pattern:
    // level | timestamp | callerFile:line | logContentJson | errorInfo
    //
    // debug/info/warn:
    // INFO | 2006-01-02 15:04:05.000 | main.go:30 | {"msg":"connected.", kvs:{}}
    //
    // error:
    // ERROR | 2006-01-02 15:04:05.000 | main.go:30 | {"msg":"connected.", kvs:{}} | rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 33.1.33.1:10246: i/o timeout"
    // [ ERROR CAUSES WITH STACK ][ rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 33.1.33.1:10246: connect: connection refused"
    // github.com/opensergo/opensergo-go/pkg/client.handleReceive.func1
    //
    //	.../opensergo-go-sdk/pkg/client/opensergo_client.go:82 ][ ERROR CAUSES WITH STACK END ]
    func AssembleMsgSeparateFormat(callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{})
```

## Samples

### sample_simple_print
If you only want to use the Logger provided by default, you can new a ConsoleLogger, like following:
``` go
    // default ConsoleLogger
    // logging.NewDefaultConsoleLogger(logging.DebugLevel)
    // custom ConsoleLogger
    logging.NewConsoleLogger(logging.DebugLevel, logging.JsonFormat, true)
    logging.Error(errors.New("errors.New"), "this is error log in printErrorStack()")
    logging.Info("this is info log in printErrorStack()")
    logging.Debug("this is debug log in printErrorStack()")
```
For detail, please refer to [sample_simple_print](./samples/sample_simple_print)

### sample_print_impl_logging
If you want to replace the logger implement where has integrated this logging component,   
you can implement the [Logger](./logger.go) by your-self, append you Logger after invoke `ClearLoggerSlice()` in Logger, 
or invoke `SetConsoleLogger(logger Logger)` to set the global ConsoleLogger.  

For detail, please refer to [sample_print_impl_logging](./samples/sample_print_impl_logging)

### sample_print_use_logging
If you want to use this logging component to replace the other has already integrated,   
you can re-implement the other Logger by invoking the function in Logger which was instanced.

For detail, please refer to [sample_print_use_logging](./samples/sample_print_use_logging)
