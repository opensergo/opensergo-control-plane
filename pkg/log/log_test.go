package log_test

import (
	"github.com/opensergo/opensergo-control-plane/pkg/log"
	"testing"
)

func Test_WithName(t *testing.T) {
	defer log.Flush()

	newLogger := log.WithName("test")
	newLogger.Info("test")
}

func Test_WithValues(t *testing.T) {
	defer log.Flush()

	newLogger := log.WithValues("key", "value")
	newLogger.Info("test")
}

func Test_Info(t *testing.T) {
	defer log.Flush()

	log.Info("test")
	log.Info("test", "key", "value")
	log.Infof("test: %s", "format")
}
