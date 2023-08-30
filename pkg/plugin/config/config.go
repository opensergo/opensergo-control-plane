package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ExternalPlugin []*Plugin `yaml:"plugin"`
}
type Plugin struct {
	Type         string `yaml:"type"`
	ExecutionDir string `yaml:"execution_dir"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error discovering plugin information: %w", err)
	}
	config := &Config{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	return config, nil
}

func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
