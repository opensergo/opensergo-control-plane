package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/pkg/errors"
)

type PluginType int

const (
	PluginTypeUnknown       PluginType = 0
	PluginTypeCompute       PluginType = 1
	PluginTypeComputeOneWay PluginType = 2
	PluginTypeComputeTwoWay PluginType = 3
)

type Plugin struct {
	PublicID    string
	PluginType  PluginType
	ScopeID     string
	Name        string //name of plugin
	Description string
	Version     uint32
	//CreateTime *timestamp.Timestamp
	//UpdateTime *timestamp.Timestamp
}

func NewPlugin(opt ...Option) *Plugin {
	opts := GetOpts(opt...)
	p := &Plugin{
		ScopeID:     Global.String(),
		Name:        opts.withName,
		Description: opts.withDescription,
	}
	return p
}

func (p *Plugin) CreatePlugin(ctx context.Context, plugintype PluginType, opt ...Option) error {
	p.PluginType = plugintype

	opts := GetOpts(opt...)

	p.PublicID = opts.withPublicID
	if p.PublicID == "" {
		var err error
		p.PublicID, err = newPublicID(p.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// newPublicID Create a globally unique publicId for a plugin
func newPublicID(name string) (string, error) {
	publicID, err := base62.Random(10)
	if err != nil {
		return "", errors.Wrap(err, "unable to generate id")
	}
	return fmt.Sprintf("%s-%s", name, publicID), nil
}
