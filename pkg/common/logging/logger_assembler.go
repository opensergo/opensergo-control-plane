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
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	defaultLogMsgBufferSize = 256
)

func caller(depth int) (file string, line int) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}

	// extract
	if osType := runtime.GOOS; osType == "windows" {
		file = strings.ReplaceAll(file, "\\", "/")
	}
	idx := strings.LastIndex(file, "/")
	file = file[idx+1:]
	return
}

// toSafeJSONString converts to valid JSON string, as the original string may contain '\\', '\n', '\r', '\t' and so on.
func toSafeJSONString(s string) []byte {
	if data, err := json.Marshal(json.RawMessage(s)); err == nil {
		return data
	} else {
		return []byte("\"" + s + "\"")
	}
}

// AssembleMsg Assemble log-msg
//
// if logFormat is not matched, return JsonFormat as default
func AssembleMsg(logFormat LogFormat, callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{}) string {
	switch logFormat {
	case SeparateFormat:
		return AssembleMsgSeparateFormat(callerDepth, logLevel, msg, err, errWithStack, keysAndValues...)
	case JsonFormat:
		return AssembleMsgJsonFormat(callerDepth, logLevel, msg, err, errWithStack, keysAndValues...)
	default:
		return AssembleMsgJsonFormat(callerDepth, logLevel, msg, err, errWithStack, keysAndValues...)
	}
}

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
func AssembleMsgJsonFormat(callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{}) string {
	sb := strings.Builder{}
	sb.Grow(defaultLogMsgBufferSize)

	file, line := caller(callerDepth)
	timeStr := time.Now().Format("2006-01-02 15:04:05.000")
	callerStr := fmt.Sprintf("%s:%d", file, line)

	sbLevel := strings.Builder{}
	sbLevel.Grow(defaultLogMsgBufferSize)
	sbLevel.WriteByte('"')
	sbLevel.WriteString("logLevel")
	sbLevel.WriteByte('"')
	sbLevel.WriteByte(':')
	sbLevel.WriteByte('"')
	sbLevel.WriteString(logLevel)
	sbLevel.WriteByte('"')

	sbTimestamp := strings.Builder{}
	sbTimestamp.Grow(defaultLogMsgBufferSize)
	sbTimestamp.WriteByte(',')
	sbTimestamp.WriteByte('"')
	sbTimestamp.WriteString("timestamp")
	sbTimestamp.WriteByte('"')
	sbTimestamp.WriteByte(':')
	sbTimestamp.WriteByte('"')
	sbTimestamp.WriteString(fmt.Sprintf("%-23s", timeStr))
	sbTimestamp.WriteByte('"')

	sbCaller := strings.Builder{}
	sbCaller.Grow(defaultLogMsgBufferSize)
	sbCaller.WriteByte(',')
	sbCaller.WriteByte('"')
	sbCaller.WriteString("caller")
	sbCaller.WriteByte('"')
	sbCaller.WriteByte(':')
	sbCaller.WriteByte('"')
	sbCaller.WriteString(callerStr)
	sbCaller.WriteByte('"')

	sbMsgJson := strings.Builder{}
	sbMsgJson.Grow(defaultLogMsgBufferSize * 2)
	sbMsgJson.WriteByte(',')
	sbMsgJson.WriteByte('"')
	sbMsgJson.WriteString("msg")
	sbMsgJson.WriteByte('"')
	sbMsgJson.WriteByte(':')
	sbMsgJson.WriteByte('"')
	sbMsgJson.WriteString(msg)
	sbMsgJson.WriteByte('"')

	sbKvs := strings.Builder{}
	sbKvs.Grow(defaultLogMsgBufferSize)
	kvLen := len(keysAndValues)
	if kvLen&1 != 0 {
		sbKvs.WriteByte(',')
		sbKvs.WriteByte('"')
		sbKvs.WriteString("kvs")
		sbKvs.WriteByte('"')
		sbKvs.WriteByte(':')
		sbKvs.Write(toSafeJSONString(fmt.Sprintf("%+v", keysAndValues)))
	} else if kvLen != 0 {
		for i := 0; i < kvLen; {
			k := keysAndValues[i]
			v := keysAndValues[i+1]
			kStr, kIsStr := k.(string)
			if !kIsStr {
				kStr = fmt.Sprintf("%+v", k)
			}
			sbKvs.WriteByte(',')
			sbKvs.Write(toSafeJSONString(kStr))
			sbKvs.WriteByte(':')
			switch v.(type) {
			case string:
				sbKvs.Write(toSafeJSONString(v.(string)))
			case error:
				sbKvs.Write(toSafeJSONString(v.(error).Error()))
			default:
				if vbs, err := json.Marshal(v); err != nil {
					sbKvs.Write(toSafeJSONString(fmt.Sprintf("%+v", v)))
				} else {
					sbKvs.Write(vbs)
				}
			}
			i = i + 2
		}
	}
	sbMsgJson.WriteString(sbKvs.String())

	sbError := strings.Builder{}
	sbError.Grow(defaultLogMsgBufferSize)
	if err != nil && err.Error() != "" {
		sbMsgJson.WriteByte(',')
		sbMsgJson.WriteByte('"')
		sbMsgJson.WriteString("error")
		sbMsgJson.WriteByte('"')
		sbMsgJson.WriteByte(':')
		sbMsgJson.WriteByte('"')
		sbMsgJson.WriteString(strings.ReplaceAll(err.Error(), "\"", "\\\""))
		sbMsgJson.WriteByte('"')

		if errWithStack {
			sbError.WriteString("[ ")
			sbError.WriteString("ERROR CAUSES WITH STACK")
			sbError.WriteString(" ]")
			sbError.WriteString("[ ")
			sbError.WriteString(fmt.Sprintf("%+v", err))
			sbError.WriteString(" ]")
			sbError.WriteString("[ ")
			sbError.WriteString("ERROR CAUSES WITH STACK END")
			sbError.WriteString(" ]")
		}
	}

	sbLog := strings.Builder{}
	sbLog.Grow(sbLevel.Len() + sbTimestamp.Len() + sbCaller.Len() + sbMsgJson.Len() + sbError.Len() + 4)
	sbLog.WriteString("{")
	sbLog.WriteString(sbLevel.String())
	sbLog.WriteString(sbTimestamp.String())
	sbLog.WriteString(sbCaller.String())
	sbLog.WriteString(sbMsgJson.String())
	sbLog.WriteString("}")

	if sbError.Len() > 0 {
		sbLog.WriteString("\n")
		sbLog.WriteString(sbError.String())
	}

	return sbLog.String()
}

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
func AssembleMsgSeparateFormat(callerDepth int, logLevel, msg string, err error, errWithStack bool, keysAndValues ...interface{}) string {

	file, line := caller(callerDepth)
	timeStr := time.Now().Format("2006-01-02 15:04:05.000")
	callerStr := fmt.Sprintf("%s:%d", file, line)

	sbLevel := strings.Builder{}
	sbLevel.Grow(defaultLogMsgBufferSize)
	sbLevel.WriteString(fmt.Sprintf("%-5s", logLevel))

	sbTimestamp := strings.Builder{}
	sbTimestamp.Grow(defaultLogMsgBufferSize)
	sbTimestamp.WriteString(" | ")
	sbTimestamp.WriteString(fmt.Sprintf("%-23s", timeStr))

	sbCaller := strings.Builder{}
	sbCaller.Grow(defaultLogMsgBufferSize)
	sbCaller.WriteString(" | ")
	sbCaller.WriteString(callerStr)

	sbLogJson := strings.Builder{}
	sbLogJson.Grow(defaultLogMsgBufferSize * 2)
	sbLogJson.WriteString(" | ")
	sbLogJson.WriteString("{")
	sbLogJson.WriteByte('"')
	sbLogJson.WriteString("msg")
	sbLogJson.WriteByte('"')
	sbLogJson.WriteByte(':')
	sbLogJson.WriteByte('"')
	sbLogJson.WriteString(msg)
	sbLogJson.WriteByte('"')

	sbKvs := strings.Builder{}
	sbKvs.Grow(defaultLogMsgBufferSize)
	kvLen := len(keysAndValues)
	if kvLen&1 != 0 {
		sbKvs.WriteByte(',')
		sbKvs.WriteByte('"')
		sbKvs.WriteString("kvs")
		sbKvs.WriteByte('"')
		sbKvs.WriteByte(':')
		sbKvs.Write(toSafeJSONString(fmt.Sprintf("%+v", keysAndValues)))
	} else if kvLen != 0 {
		for i := 0; i < kvLen; {
			k := keysAndValues[i]
			v := keysAndValues[i+1]
			kStr, kIsStr := k.(string)
			if !kIsStr {
				kStr = fmt.Sprintf("%+v", k)
			}
			sbKvs.WriteByte(',')
			sbKvs.Write(toSafeJSONString(kStr))
			sbKvs.WriteByte(':')
			switch v.(type) {
			case string:
				sbKvs.Write(toSafeJSONString(v.(string)))
			case error:
				sbKvs.Write(toSafeJSONString(v.(error).Error()))
			default:
				if vbs, err := json.Marshal(v); err != nil {
					sbKvs.Write(toSafeJSONString(fmt.Sprintf("%+v", v)))
				} else {
					sbKvs.Write(vbs)
				}
			}
			i = i + 2
		}
	}

	sbLogJson.WriteString(sbKvs.String())
	sbLogJson.WriteString("}")

	sbError := strings.Builder{}
	sbError.Grow(defaultLogMsgBufferSize)
	if err != nil && err.Error() != "" {
		sbError.WriteString(" | ")
		sbError.WriteString(err.Error())
		sbError.WriteString("\n")

		if errWithStack {
			sbError.WriteString("[ ")
			sbError.WriteString("ERROR CAUSES WITH STACK")
			sbError.WriteString(" ]")
			sbError.WriteString("[ ")
			sbError.WriteString(fmt.Sprintf("%+v", err))
			sbError.WriteString(" ]")
			sbError.WriteString("[ ")
			sbError.WriteString("ERROR CAUSES WITH STACK END")
			sbError.WriteString(" ]")
		}
	}

	sbLog := strings.Builder{}
	sbLog.Grow(sbLevel.Len() + sbTimestamp.Len() + sbCaller.Len() + sbLogJson.Len() + sbError.Len() + 4)
	sbLog.WriteString(sbLevel.String())
	sbLog.WriteString(sbTimestamp.String())
	sbLog.WriteString(sbCaller.String())
	sbLog.WriteString(sbLogJson.String())

	if sbError.Len() > 0 {
		sbLog.WriteString(sbError.String())
	}

	sbLog.WriteString("\n")
	return sbLog.String()
}
