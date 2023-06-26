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

package controller

import (
	"strings"

	"github.com/alibaba/sentinel-golang/logging"
	"github.com/go-logr/logr"
)

// noopLogSink is a logr.Logger that's always disabled, and does nothing.
type noopLogSink struct{}

func (n noopLogSink) Init(info logr.RuntimeInfo) {
	//TODO implement me
}

func (n noopLogSink) Enabled(level int) bool {
	//TODO implement me
	panic("implement me")
}

func (n noopLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (n noopLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (n noopLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return n
}

func (n noopLogSink) WithName(name string) logr.LogSink {
	return n
}

var disabledLogSink = noopLogSink{}

type k8SLogSink struct {
	l             logging.Logger
	level         logging.Level
	names         []string
	keysAndValues []interface{}
}

func (k k8SLogSink) Init(info logr.RuntimeInfo) {
	// TODO: init the k8s logger using info
}

func (k k8SLogSink) Enabled(level int) bool {
	return true
}

func (k k8SLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, k.keysAndValues...)
	switch k.level {
	case logging.WarnLevel:
		k.l.Warn(k.buildNames()+msg, keysAndValues...)
	case logging.InfoLevel:
		k.l.Info(k.buildNames()+msg, keysAndValues...)
	case logging.DebugLevel:
		k.l.Debug(k.buildNames()+msg, keysAndValues...)
	default:
		k.l.Info(k.buildNames()+msg, keysAndValues...)
	}
}

func (k *k8SLogSink) buildNames() string {
	size := len(k.names)
	if size == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, name := range k.names {
		sb.WriteString(name)
		if i == size-1 {
			continue
		}
		sb.WriteString(".")
	}
	sb.WriteString(" ")
	return sb.String()
}

func (k *k8SLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, k.keysAndValues...)
	k.l.Error(err, k.buildNames()+msg, keysAndValues...)
}

func (k *k8SLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	names := make([]string, len(k.names))
	copy(names, k.names)
	kvs := make([]interface{}, len(k.keysAndValues))
	copy(kvs, k.keysAndValues)
	kvs = append(kvs, keysAndValues...)
	return &k8SLogSink{
		l:             k.l,
		level:         k.level,
		names:         names,
		keysAndValues: kvs,
	}
}

func (k *k8SLogSink) WithName(name string) logr.LogSink {
	names := make([]string, len(k.names))
	copy(names, k.names)
	names = append(names, name)
	kvs := make([]interface{}, len(k.keysAndValues))
	copy(kvs, k.keysAndValues)
	return &k8SLogSink{
		l:             k.l,
		level:         k.level,
		names:         names,
		keysAndValues: kvs,
	}
}

func (k *k8SLogSink) Logger() logr.Logger {
	return logr.Logger{}.WithSink(k).WithName(k.buildNames())
}
