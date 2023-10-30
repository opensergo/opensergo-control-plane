package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-secure-stdlib/pluginutil/v2"
)

func CreatePlugin(ctx context.Context, pluginName, pluginSetName string, opt ...Option) (raw any, cleanup func() error, err error) {
	raw, cleanup, err = createPlugin(ctx, pluginName, pluginSetName, opt...)
	if err != nil {
		return nil, cleanup, err
	}

	//var ok bool
	//hp, ok := raw.(pb.StreamGreeterClient)
	//if !ok {
	//	return nil, cleanup, fmt.Errorf("error converting rpc plugin of type %T to normal wrapper", raw)
	//}

	return raw, cleanup, nil
}

func createPlugin(
	ctx context.Context,
	pluginName string,
	pluginSetName string,
	opt ...Option,
) (
	raw any,
	cleanup func() error,
	retErr error,
) {
	defer func() {
		if retErr != nil && cleanup != nil {
			_ = cleanup()
		}
	}()

	pluginName = strings.ToLower(pluginName)

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing plugin options: %w", err)
	}

	// First, scan available plugins, then find the right one to use
	pluginMap, err := pluginutil.BuildPluginMap(
		append(
			opts.withPluginOptions,
			pluginutil.WithPluginClientCreationFunc(
				func(pluginPath string, _ ...pluginutil.Option) (*plugin.Client, error) {
					return NewPluginClient(pluginPath, pluginSetName, WithLogger(opts.withLogger))
				}),
		)...)
	if err != nil {
		return nil, nil, fmt.Errorf("error building plugin map: %w", err)
	}

	// Create the plugin and cleanup func
	plugClient, cleanup, err := pluginutil.CreatePlugin(pluginMap[pluginName], opts.withPluginOptions...)
	if err != nil {
		return nil, cleanup, err
	}

	switch client := plugClient.(type) {
	case plugin.ClientProtocol:
		raw, err = client.Dispense(pluginSetName)
		if err != nil {
			return nil, cleanup, fmt.Errorf("error dispensing %q plugin: %w", pluginSetName, err)
		}
	default:
		return nil, cleanup, fmt.Errorf("unable to understand type %T of raw plugin", raw)
	}

	return raw, cleanup, nil
}
