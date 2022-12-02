package log_test

import (
	"fmt"
	"github.com/opensergo/opensergo-control-plane/pkg/log"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Options_AddFlags(t *testing.T) {
	opts := log.NewOptions()
	fs := pflag.NewFlagSet("test", pflag.ExitOnError)
	opts.AddFlags(fs)

	args := []string{"--log.level=debug"}
	err := fs.Parse(args)
	assert.Nil(t, err)

	assert.Equal(t, "debug", opts.Level)
}

func Test_Options_Validate(t *testing.T) {
	opts := log.NewOptions(log.WithFormat("test"), log.WithLevel("test"))
	errs := opts.Validate()
	expected := `[unrecognized level: "test" unrecognized format: "test"]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))

	opts = log.NewOptions()
	errs = opts.Validate()
	assert.Equal(t, []error(nil), errs)
}
