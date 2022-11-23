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

package main

import (
	"github.com/opensergo/opensergo-control-plane/pkg/common/logging"
	"log"
	"os"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	DebugEnabled() bool

	// Info logs a non-error message with the given key/value pairs as context.
	//
	// The msg argument should be used to add some constant description to
	// the log line.  The key/value pairs can then be used to add additional
	// variable information.  The key/value pairs should alternate string
	// keys and arbitrary values.
	Info(msg string, keysAndValues ...interface{})
	InfoEnabled() bool

	Warn(msg string, keysAndValues ...interface{})
	WarnEnabled() bool

	Error(err error, msg string, keysAndValues ...interface{})
	ErrorEnabled() bool
}

type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type SentinelLogger struct {
	l           *log.Logger
	loggerLevel Level
}

func (sl *SentinelLogger) Debug(msg string, keysAndValues ...interface{}) {
	sl.l.Print(msg)
}

func (sl *SentinelLogger) DebugEnabled() bool {
	return DebugLevel >= sl.loggerLevel
}

func (sl *SentinelLogger) Info(msg string, keysAndValues ...interface{}) {
	sl.l.Print(msg)

}
func (sl *SentinelLogger) InfoEnabled() bool {
	return InfoLevel >= sl.loggerLevel
}

func (sl *SentinelLogger) Warn(msg string, keysAndValues ...interface{}) {
	sl.l.Print(msg)
}

func (sl *SentinelLogger) WarnEnabled() bool {
	return WarnLevel >= sl.loggerLevel
}

func (sl *SentinelLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	sl.l.Print(msg)
}

func (sl *SentinelLogger) ErrorEnabled() bool {
	return ErrorLevel >= sl.loggerLevel
}

func NewDefaultSentinelLogger() *SentinelLogger {
	return &SentinelLogger{
		l:           log.New(os.Stdout, "", 0),
		loggerLevel: DebugLevel,
	}
}

type LoggerAdaptor struct {
	l              *SentinelLogger
	loggerLevel    logging.Level
	loggerFormat   logging.LogFormat
	errorWithStack bool
}

func (la *LoggerAdaptor) Print(msg string) {
	logger := log.New(os.Stdout, "", 0)
	logger.Print(msg)
}

func (la *LoggerAdaptor) Debug(msg string, keysAndValues ...interface{}) {
	la.l.Debug(msg, keysAndValues)
}

func (la *LoggerAdaptor) DebugEnabled() bool {
	return la.l.DebugEnabled()
}

func (la *LoggerAdaptor) Info(msg string, keysAndValues ...interface{}) {
	la.l.Info(msg, keysAndValues)
}

func (la *LoggerAdaptor) InfoEnabled() bool {
	return la.l.InfoEnabled()
}

func (la *LoggerAdaptor) Warn(msg string, keysAndValues ...interface{}) {
	la.l.Warn(msg, keysAndValues)
}

func (la *LoggerAdaptor) WarnEnabled() bool {
	return la.l.WarnEnabled()
}

func (la *LoggerAdaptor) Error(err error, msg string, keysAndValues ...interface{}) {
	la.l.Error(err, msg, keysAndValues)
}

func (la *LoggerAdaptor) ErrorEnabled() bool {
	return la.l.ErrorEnabled()
}

func NewLoggerAdaptor() *LoggerAdaptor {
	return &LoggerAdaptor{
		l:              NewDefaultSentinelLogger(),
		loggerLevel:    logging.DebugLevel,
		loggerFormat:   logging.JsonFormat,
		errorWithStack: true,
	}
}
