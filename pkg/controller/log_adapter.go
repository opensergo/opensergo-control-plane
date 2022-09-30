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

// noopLogger is a logr.Logger that's always disabled, and does nothing.
type noopLogger struct{}

func (l *noopLogger) Enabled() bool                             { return false }
func (l *noopLogger) Info(_ string, _ ...interface{})           {}
func (l *noopLogger) Error(_ error, _ string, _ ...interface{}) {}
func (l *noopLogger) V(_ int) logr.Logger                       { return l }
func (l *noopLogger) WithValues(_ ...interface{}) logr.Logger   { return l }
func (l *noopLogger) WithName(_ string) logr.Logger             { return l }

var disabledLogger = &noopLogger{}

type k8SLogger struct {
	l             logging.Logger
	level         logging.Level
	names         []string
	keysAndValues []interface{}
}

func (k *k8SLogger) buildNames() string {
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

func (k *k8SLogger) Info(msg string, keysAndValues ...interface{}) {
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

func (k *k8SLogger) Enabled() bool {
	return true
}

func (k *k8SLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, k.keysAndValues...)
	k.l.Error(err, k.buildNames()+msg, keysAndValues...)
}

func (k *k8SLogger) V(level int) logr.Logger {
	if k.Enabled() {
		names := make([]string, len(k.names))
		copy(names, k.names)
		kvs := make([]interface{}, len(k.keysAndValues))
		copy(kvs, k.keysAndValues)
		return &k8SLogger{
			l:             k.l,
			level:         logging.Level(level),
			names:         names,
			keysAndValues: kvs,
		}
	}
	return disabledLogger
}

func (k *k8SLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	names := make([]string, len(k.names))
	copy(names, k.names)
	kvs := make([]interface{}, len(k.keysAndValues))
	copy(kvs, k.keysAndValues)
	kvs = append(kvs, keysAndValues...)
	return &k8SLogger{
		l:             k.l,
		level:         k.level,
		names:         names,
		keysAndValues: kvs,
	}
}

func (k *k8SLogger) WithName(name string) logr.Logger {
	names := make([]string, len(k.names))
	copy(names, k.names)
	names = append(names, name)
	kvs := make([]interface{}, len(k.keysAndValues))
	copy(kvs, k.keysAndValues)
	return &k8SLogger{
		l:             k.l,
		level:         k.level,
		names:         names,
		keysAndValues: kvs,
	}
}
