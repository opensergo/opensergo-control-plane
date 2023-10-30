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
	b := NewStreamPlugin()
	if err := plugin.ServePlugin(b); err != nil {
		fmt.Println("Error serving plugin", err)
		os.Exit(1)
	}
	os.Exit(0)
}

var (
	_ stream_plugin.Stream = (*stream_plugin.StreamPluginServer)(nil)
)

type BuiltinPlugin struct {
	*stream_plugin.StreamPluginServer
}

func NewStreamPlugin() *BuiltinPlugin {
	return &BuiltinPlugin{
		StreamPluginServer: &stream_plugin.StreamPluginServer{},
	}
}
