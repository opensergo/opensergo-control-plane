package plugin

import (
	"errors"
	"os/exec"

	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin"
	stream_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/stream"

	"github.com/hashicorp/go-plugin"
)

// HandshakeConfig is a shared config that can be used regardless of plugin, to
// avoid having to know type-specific things about each plugin
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "OPENSERGO_STREAM_PLUGIN",
	MagicCookieValue: "opensergo-plugin",
}

// ServePlugin starts a plugin server
func ServePlugin(svc any, opt ...Option) error {
	opts, err := getOpts(opt...)
	if err != nil {
		return err
	}

	plugins := make(map[string]plugin.Plugin)
	if streamSvc, ok := svc.(stream_plugin.Stream); ok {
		streamServiceServer, err := stream_plugin.NewStreamPluginServiceServer(streamSvc)
		if err != nil {
			return err
		}
		plugins[builtin.StreamServicePluginSetName] = streamServiceServer
	}

	if len(plugins) == 0 {
		return errors.New("no valid plugin server provided")
	}

	//opts.withLogger.Info("Info NewPluginServer")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: plugins,
		},
		Logger:     opts.withLogger,
		GRPCServer: plugin.DefaultGRPCServer,
	})
	return nil
}

func NewPluginClient(pluginPath string, setName string, opt ...Option) (*plugin.Client, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	var set plugin.PluginSet

	switch setName {
	case builtin.StreamServicePluginSetName:
		set = plugin.PluginSet{builtin.StreamServicePluginSetName: &stream_plugin.StreamPlugin{}}
	}

	//fmt.Println("NewPluginClient")
	//opts.withLogger.Info("Info NewPluginClient")
	//opts.withLogger.Debug("Debug NewPluginClient")

	return plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: set,
		},
		Cmd: exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		Logger:   opts.withLogger,
		AutoMTLS: true,
	}), nil
}
