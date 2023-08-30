package plugin

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-secure-stdlib/pluginutil/v2"

	"github.com/hashicorp/go-hclog"
)

// getOpts iterates the inbound Options and returns a struct
func getOpts(opt ...Option) (*options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(opts); err != nil {
			return nil, fmt.Errorf("error running option function: %w", err)
		}
	}
	return opts, nil
}

// Option - a type that wraps an interface for compile-time safety but can
// contain an option for this package or for wrappers implementing this
// interface.
type Option func(*options) error

type options struct {
	withPluginOptions []pluginutil.Option
	withLogger        hclog.Logger
}

func getDefaultOptions() *options {
	return &options{}
}

// WithPluginOptions allows providing plugin-related (as opposed to
// configutil-related) options
func WithPluginOptions(opts ...pluginutil.Option) Option {
	return func(o *options) error {
		o.withPluginOptions = append(o.withPluginOptions, opts...)
		return nil
	}
}

// WithLogger allows passing a logger to the plugin library for debugging
func WithLogger(logger hclog.Logger) Option {
	return func(o *options) error {
		o.withLogger = logger
		return nil
	}
}
func NewLogger(name string) hclog.Logger {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   name,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	return logger
}
