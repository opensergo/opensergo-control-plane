package main

import (
	"fmt"
	"os"

	ratelimit_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/ratelimit"

	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/plugin"
)

func main() {
	//log := hclog.New(&hclog.LoggerOptions{
	//	Output:     os.Stderr,
	//	Level:      hclog.Trace,
	//	JSONFormat: true,
	//}), plugin.WithLogger(log)
	b := NewBuiltinPlugin()
	if err := plugin.ServePlugin(b); err != nil {
		fmt.Println("Error serving plugin", err)
		os.Exit(1)
	}
	os.Exit(0)
}

var (
	_ ratelimit_plugin.RateLimit = (*ratelimit_plugin.RateLimitPluginServer)(nil)
)

type BuiltinPlugin struct {
	*ratelimit_plugin.RateLimitPluginServer
}

func NewBuiltinPlugin() *BuiltinPlugin {
	return &BuiltinPlugin{
		RateLimitPluginServer: &ratelimit_plugin.RateLimitPluginServer{},
	}
}
