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

// when application has use logging-component in opensergo-go-common
// and import the other logging-component such as sentinel logger
// how to print logs with the sentinel logger through logging-component in opensergo-go-common
func main() {
	logger, _ := logging.NewDefaultFileLogger(logging.InfoLevel)
	logging.AppendLoggerSlice(logger)

	adaptor := NewLoggerAdaptor()
	logging.ClearLoggerSlice()
	logging.AppendLoggerSlice(adaptor)
	logging.Error(errors.New("errors.New"), "this is error log implement logging in sentinelLoggerAdaptor")
	logging.Info("this is info log implement logging in sentinelLoggerAdaptor")
	logging.Debug("this is debug log implement logging in sentinelLoggerAdaptor")
}
