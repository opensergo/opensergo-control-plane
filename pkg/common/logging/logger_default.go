// Copyright 2022, OpenSergo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
)

// Level represents the level of logging.
type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type LogFormat uint8

const (
	JsonFormat LogFormat = iota
	SeparateFormat
)

var (
	ConsoleLogFormat = JsonFormat
	ConsoleLogLevel  = InfoLevel
)

const (
	DefaultLogCallerDepth = 8
	DefaultErrorWithStack = false
)

// DefaultLogger default logger implements Logger
type DefaultLogger struct {
	Logger         *log.Logger
	LoggerLevel    Level
	LoggerFormat   LogFormat
	ErrorWithStack bool
}

func (l *DefaultLogger) Print(msg string) {
	l.Logger.Println(msg)
}

func (l *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	DebugWithCallerDepth(l, l.LoggerFormat, DefaultLogCallerDepth+1, msg, keysAndValues...)
}

func (l *DefaultLogger) DebugEnabled() bool {
	return DebugLevel >= l.LoggerLevel
}

func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	InfoWithCallerDepth(l, l.LoggerFormat, DefaultLogCallerDepth+1, msg, keysAndValues...)
}

func (l *DefaultLogger) InfoEnabled() bool {
	return InfoLevel >= l.LoggerLevel
}

func (l *DefaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	WarnWithCallerDepth(l, l.LoggerFormat, DefaultLogCallerDepth+1, msg, keysAndValues...)
}

func (l *DefaultLogger) WarnEnabled() bool {
	return WarnLevel >= l.LoggerLevel
}

func (l *DefaultLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	ErrorWithCallerDepth(l, l.LoggerFormat, DefaultLogCallerDepth+1, err, l.ErrorWithStack, msg, keysAndValues...)
}

func (l *DefaultLogger) ErrorEnabled() bool {
	return ErrorLevel >= l.LoggerLevel
}

// NewDefaultFileLogger
// LoggerLevel DefaultLogLevel
func NewDefaultFileLogger(logLevel Level) (Logger, error) {
	// log dir of 'OpenSergo universal transport service'
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		userHomeDir = ""
	}
	filePath := filepath.Join(userHomeDir, "logs", "opensergo", "opensergo-universal-transport-service.log")
	mkdirLogFile(filePath)

	return NewFileLogger(filePath, logLevel, JsonFormat, DefaultErrorWithStack)
}

// create a file to write log.
func mkdirLogFile(filePath string) {
	dir := filepath.Dir(filePath)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				Error(err, dir+" create directory error!")
			}
		}
	}
	if fileInfo, _ := os.Stat(dir); !fileInfo.IsDir() {
		Error(errors.New(dir+" is not a directiry!"), dir+" is not a directiry!")
	}

}

// NewFileLogger new a Logger which write logs into file. And append the LoggerSlice.
// filepath is the full path(absolute path). eg: /root/logs/opensergo/opensergo-universal-transport-service.log
func NewFileLogger(filepath string, loggerLevel Level, logFormat LogFormat, errorWithStack bool) (Logger, error) {
	mkdirLogFile(filepath)
	logFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defaultLogger := &DefaultLogger{
		Logger:         log.New(logFile, "", 0),
		LoggerLevel:    loggerLevel,
		LoggerFormat:   logFormat,
		ErrorWithStack: errorWithStack,
	}
	AppendLoggerSlice(defaultLogger)
	return defaultLogger, err
}

// NewDefaultConsoleLogger new a default ConsoleLogger to print logs in console.
// And there is only one ConsoleLogger instance in the Global.
func NewDefaultConsoleLogger(logLevel Level) Logger {
	defaultLogger := NewConsoleLogger(logLevel, ConsoleLogFormat, DefaultErrorWithStack)
	SetConsoleLogger(defaultLogger)
	return defaultLogger
}

// NewConsoleLogger new a ConsoleLogger to print logs in console.
// And there is only one ConsoleLogger instance in the Global.
func NewConsoleLogger(logLevel Level, logFormat LogFormat, errorWithStack bool) Logger {
	defaultLogger := &DefaultLogger{
		Logger:         log.New(os.Stdout, "", 0),
		LoggerLevel:    logLevel,
		LoggerFormat:   logFormat,
		ErrorWithStack: errorWithStack,
	}
	SetConsoleLogger(defaultLogger)
	return defaultLogger
}
