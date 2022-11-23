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
	"github.com/pkg/errors"
)

// when application has use logging-component Self-contained
// and import logging-component in opensergo-go-common at the same time
// how to print logs with the format from logging-component in opensergo-go-common through logging-component Self-contained
func main() {
	logger := logging.NewConsoleLogger(logging.DebugLevel, logging.JsonFormat, true)
	sentinelLogger := NewSentinelLogger(logger, logging.JsonFormat, true)
	sentinelLogger.Error(errors.New("errors.New"), "this is error log use logging in sentinelLogger")
	sentinelLogger.Info("this is info log use logging in sentinelLogger")
	sentinelLogger.Debug("this is debug log use logging in sentinelLogger")
}
