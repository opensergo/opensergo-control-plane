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

package util

import (
	"strconv"
	"strings"

	pb "github.com/opensergo/opensergo-control-plane/pkg/proto/fault_tolerance/v1"
)

func Str2MillSeconds(timeStr string) (int64, error) {
	if strings.Contains(timeStr, "ms") && len(timeStr) > 2 {
		return strconv.ParseInt(timeStr[:len(timeStr)-2], 10, 64)
	} else if strings.Contains(timeStr, "s") && len(timeStr) > 1 {
		millSec, err := strconv.ParseInt(timeStr[:len(timeStr)-1], 10, 64)
		return millSec * 10, err
	} else if strings.Contains(timeStr, "minute") && len(timeStr) > 6 {
		millSec, err := strconv.ParseInt(timeStr[:len(timeStr)-6], 10, 64)
		return millSec * 60 * 10, err
	} else if strings.Contains(timeStr, "min") && len(timeStr) > 3 {
		millSec, err := strconv.ParseInt(timeStr[:len(timeStr)-3], 10, 64)
		return millSec * 60 * 10, err
	} else if strings.Contains(timeStr, "m") && len(timeStr) > 1 {
		millSec, err := strconv.ParseInt(timeStr[:len(timeStr)-1], 10, 64)
		return millSec * 60 * 10, err
	} else if strings.Contains(timeStr, "h") && len(timeStr) > 1 {
		millSec, err := strconv.ParseInt(timeStr[:len(timeStr)-1], 10, 64)
		return millSec * 60 * 60 * 10, err
	}
	return strconv.ParseInt(timeStr, 10, 64)
}

func Str2LimitNode(lnStr string) pb.ConcurrencyLimitStrategy_LimitMode {
	lnStr = strings.ToUpper(lnStr)
	switch lnStr {
	case "LOCAL":
		return pb.ConcurrencyLimitStrategy_MODE_LOCAL
	case "GLOBAL":
		return pb.ConcurrencyLimitStrategy_MODE_GLOBAL
	default:
		return pb.ConcurrencyLimitStrategy_MODE_UNKNOWN
	}
}

func Str2CBStrategy(cbStr string) pb.CircuitBreakerStrategy_Strategy {
	cbStr = strings.ToUpper(cbStr)
	switch cbStr {
	case strings.ToUpper("SlowRequestRatio"):
		return pb.CircuitBreakerStrategy_STRATEGY_SLOW_REQUEST_RATIO
	case strings.ToUpper("ErrorRequestRatio"):
		return pb.CircuitBreakerStrategy_STRATEGY_ERROR_REQUEST_RATIO
	default:
		return pb.CircuitBreakerStrategy_STRATEGY_UNKNOWN
	}
}

func RatioStr2Float(ratioStr string) (float64, error) {
	if strings.Contains(ratioStr, "%") && len(ratioStr) > 1 {
		num, err := strconv.ParseFloat(ratioStr[:len(ratioStr)-1], 64)
		return num / 100.0, err
	}
	return strconv.ParseFloat(ratioStr, 64)
}
