package options

import (
	"github.com/spf13/pflag"
	"log"
	"os"
)

type Options struct {
	ConfigPushMaxAttempt int
	flags                *pflag.FlagSet
}

func (o *Options) AddFlag() {
	o.flags.IntVar(&o.ConfigPushMaxAttempt, "ConfigPushMaxAttempt", 3, "max times for pushing config after timeout error")
}

func (o *Options) Parse() error {
	err := o.flags.Parse(os.Args)
	return err
}

func NewOption() (*Options, error) {
	o := &Options{
		flags: pflag.NewFlagSet("sergo-flag", pflag.ExitOnError),
	}
	o.AddFlag()
	err := o.Parse()
	if err != nil {
		log.Fatalf("Parse flag failure: %s", err)
	}
	return o, err
}
