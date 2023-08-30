package main

import (
	"fmt"
	"os"

	stream_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/stream"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/plugin"
)

func main() {
	//log := hclog.New(&hclog.LoggerOptions{
	//	Output:     os.Stderr,
	//	Level:      hclog.Trace,
	//	JSONFormat: true,
	//}), plugin.WithLogger(log)
	streamplugin := &stream_plugin.StreamPluginServer{}
	if err := plugin.ServePlugin(streamplugin); err != nil {
		fmt.Println("Error serving plugin", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// 为了加载多个plugin，现在不用这个，后面plugin多了
// 就加在StreamPlugin中，然后使用StreamPlugin传入ServePlugin里
type StreamPlugin struct {
	*stream_plugin.StreamPluginServer
}

func NewStreamPlugin() *StreamPlugin {
	return &StreamPlugin{
		&stream_plugin.StreamPluginServer{},
	}
}
