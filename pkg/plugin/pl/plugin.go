package pl

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/opensergo/opensergo-control-plane/pkg/plugin"
	pluginclient "github.com/opensergo/opensergo-control-plane/pkg/plugin/client"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/config"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin"

	plugin2 "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/plugin"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/plugin/store"

	"github.com/hashicorp/go-secure-stdlib/pluginutil/v2"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/util"
)

type PluginServer struct {
	Context       context.Context
	ContextCancel context.CancelFunc
	Config        *config.Config
	EnabledPlugin []*EnabledPlugin
	Client        *pluginclient.PluginClientRegistry
	ShutdownCh    chan struct{}
}

type EnabledPlugin struct {
	*store.Plugin
	PluginName    string
	ShutdownFuncs []func() error
}

func NewPluginServer() *PluginServer {
	ctx, cancel := context.WithCancel(context.Background())
	ps := &PluginServer{
		Context:       ctx,
		ContextCancel: cancel,
		EnabledPlugin: make([]*EnabledPlugin, 0),
		Client:        &pluginclient.PluginClientRegistry{},
		ShutdownCh:    MakeShutdownCh(),
	}
	//go func() {
	//	<-ps.ShutdownCh
	//	cancel()
	//	if err := ps.RunShutdownFuncs(); err != nil {
	//		log.Fatalln("Error:", err.Error())
	//	}
	//	log.Println("Server shutdown")
	//}()
	return ps
}

func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})
	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			<-shutdownCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}

func (p *PluginServer) RunShutdownFuncs() error {
	for i := range p.EnabledPlugin {
		for _, f := range p.EnabledPlugin[i].ShutdownFuncs {
			if err := f(); err != nil {
				return fmt.Errorf("error shutting down plugin %s: %w", p.EnabledPlugin[i].PluginName, err)
			}
			p.Client.DeletePluginClient(p.EnabledPlugin[i].PublicID)
		}
	}
	return nil
}

func (p *PluginServer) InitPlugin() error {
	path := config.GetCurrentAbPathByCaller()
	c, err := config.ReadConfig(filepath.Join(filepath.Dir(path), "config/config.yaml"))
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	p.Config = c
	//for i, _ := range c.ExternalPlugin {
	//	p.EnabledPlugin = append(p.EnabledPlugin, EnabledPlugin{
	//		PluginName: c.ExternalPlugin[i].Type,
	//	})
	//}

	// enadle builtin plugin
	p.EnabledPlugin = append(p.EnabledPlugin, &EnabledPlugin{
		PluginName: builtin.StreamServicePluginName,
	})
	err = p.CreatePlugin()
	if err != nil {
		return fmt.Errorf("error creating plugin: %w", err)
	}

	return nil
}

type creatOption struct {
	pluginSetName string
	pluginType    store.PluginType
	executionDir  string
}

func (p *PluginServer) CreatePlugin() error {
	for _, enabledPlugin := range p.EnabledPlugin {
		switch enabledPlugin.PluginName {
		case builtin.StreamServicePluginName:
			co := &creatOption{
				pluginSetName: builtin.StreamServicePluginSetName,
				pluginType:    store.PluginTypeCompute,
				executionDir:  "",
			}
			err := p.createplugin(enabledPlugin, co)
			if err != nil {
				return fmt.Errorf("error CreatePlugin: %w", err)
			}
		default:
			fmt.Printf("unknow plugin: %s\n", enabledPlugin.PluginName)
		}
	}
	return nil
}

func (p *PluginServer) createplugin(e *EnabledPlugin, c *creatOption) error {
	logger := plugin2.NewLogger(c.pluginSetName)
	pluginName := strings.ToLower(e.PluginName)
	client, cleanup, err := plugin2.CreatePlugin(
		p.Context,
		pluginName,
		c.pluginSetName,
		plugin2.WithPluginOptions(
			pluginutil.WithPluginExecutionDirectory(c.executionDir),
			pluginutil.WithPluginsFilesystem(util.PluginPrefix, plugin.FileSystem()),
		),
		plugin2.WithLogger(logger),
	)
	e.ShutdownFuncs = append(e.ShutdownFuncs, cleanup)
	if err != nil {
		return fmt.Errorf("error creating %s plugin: %w", pluginName, err)
	}

	plg, err := p.registerPlugin(p.Context, pluginName, client, c.pluginType, store.WithDescription(fmt.Sprintf("Built-in %s plugin", pluginName)))
	if err != nil {
		return fmt.Errorf("error registering %s plugin: %w", pluginName, err)
	}
	e.Plugin = plg
	return nil
}

func (p *PluginServer) registerPlugin(ctx context.Context, name string, client interface{}, flag store.PluginType, opt ...store.Option) (*store.Plugin, error) {
	opt = append(opt, store.WithName(name))
	plg := store.NewPlugin(opt...)
	err := plg.CreatePlugin(ctx, flag, opt...)
	if err != nil {
		return nil, err
	}
	if err = p.Client.RegisterClient(plg.PublicID, client); err != nil {
		return nil, err
	}
	return plg, nil

}

func (p *PluginServer) GetPluginClient(name string) (interface{}, error) {
	//keys := make([]string, 0, len(p.Client))
	client := p.Client.RangePluginClient(name)
	if client == nil {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return client, nil
}
