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

import "github.com/opensergo/opensergo-control-plane/pkg/common/logging"

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

type SentinelLogger struct {
	l              logging.Logger
	logFormat      logging.LogFormat
	errorWithStack bool
}

func (sl *SentinelLogger) Debug(msg string, keysAndValues ...interface{}) {
	logging.DebugWithCallerDepth(sl.l, sl.logFormat, 5, msg, keysAndValues)
}

func (sl *SentinelLogger) DebugEnabled() bool {
	return sl.l.DebugEnabled()
}

func (sl *SentinelLogger) Info(msg string, keysAndValues ...interface{}) {
	logging.InfoWithCallerDepth(sl.l, sl.logFormat, 5, msg, keysAndValues)

}
func (sl *SentinelLogger) InfoEnabled() bool {
	return sl.l.InfoEnabled()
}

func (sl *SentinelLogger) Warn(msg string, keysAndValues ...interface{}) {
	logging.WarnWithCallerDepth(sl.l, sl.logFormat, 5, msg, keysAndValues)
}

func (sl *SentinelLogger) WarnEnabled() bool {
	return sl.l.WarnEnabled()
}

func (sl *SentinelLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logging.ErrorWithCallerDepth(sl.l, sl.logFormat, 5, err, sl.errorWithStack, msg, keysAndValues)
}

func (sl *SentinelLogger) ErrorEnabled() bool {
	return sl.l.ErrorEnabled()
}

func NewSentinelLogger(logger logging.Logger, logFormat logging.LogFormat, errorWithStack bool) *SentinelLogger {
	return &SentinelLogger{
		l:              logger,
		logFormat:      logFormat,
		errorWithStack: errorWithStack,
	}
}
