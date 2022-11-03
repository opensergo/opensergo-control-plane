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

// use default logger to print log
func main() {
	printErrorStack()

	logging.Print("--------------------")

	printDefault()
}

func printErrorStack() {
	logging.NewConsoleLogger(logging.DebugLevel, logging.JsonFormat, true)
	logging.Error(errors.New("errors.New"), "this is error log in printErrorStack()")
	logging.Info("this is info log in printErrorStack()")
	logging.Debug("this is debug log in printErrorStack()")
}

func printDefault() {
	logging.NewDefaultConsoleLogger(logging.DebugLevel)
	logging.Error(errors.New("errors.New"), "this is error log in printDefault()")
	logging.Info("this is info log in printDefault()")
	logging.Debug("this is debug log in printDefault()")
}
