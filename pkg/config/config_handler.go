// Copyright 2022, OpenSergo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"gopkg.in/yaml.v2"
)

// LoadConfig for start OpenSergo Server.
// Priority of loading config: LoadOption > SystemEnv > ConfigFile > DefaultConfig.
// 1. NewDefaultConfig() to get the default config.
// 2. Get ConfPath and override config from config file by overrideFromYml() .
// 3. Override config from SystemEnv by overrideFromSystemEnv().
// 4. Override config from Option by overrideFromOpts().
func LoadConfig(opts ...Option) (*OpenSergoConfig, error) {
	// default config
	mergedConfig := NewDefaultConfig()

	// update initial config from File from Opts
	tmpConfig := NewDefaultConfig()
	tmpConfig.overrideFromOpts(opts...)
	mergedConfig.ConfPath = tmpConfig.ConfPath
	if err := mergedConfig.overrideFromYml(mergedConfig.ConfPath); err != nil {
		log.Println("read config from ConfPath[{}] error", tmpConfig.ConfPath)
	}

	// update initial config from System Env
	if err := mergedConfig.overrideFromSystemEnv(); err != nil {
		return nil, err
	}

	// update initial config from func args
	mergedConfig.overrideFromOpts(opts...)

	return mergedConfig, nil
}

func (c *OpenSergoConfig) overrideFromOpts(opts ...Option) {
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(c)
		}
	}
}

func (c *OpenSergoConfig) InitOptsFromCommand() (opts []Option) {
	flag.StringVar(&c.ConfPath, "c", DefaultConfPath, "file path of config")
	flag.UintVar(&c.Port, "p", DefaultPort, "endpoint port of OpenSergo Control Plane.[ SystemEnvName: "+EnvKeyEndpointPort+" ][ YamlPath:endpointPort ]")
	flag.Parse()
	if c.ConfPath != DefaultConfPath {
		opts = append(opts, WithConfPath(c.ConfPath))
	}
	if c.Port != DefaultPort {
		opts = append(opts, WithEndpointPort(c.Port))
	}
	return opts
}

func (c *OpenSergoConfig) overrideFromYml(confPath string) error {
	_, err := os.Stat(confPath)
	if err != nil && !os.IsExist(err) {
		return err
	}
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	source := NewDefaultConfig()
	source.ConfPath = confPath
	err = yaml.Unmarshal(content, source)
	if err != nil {
		return err
	}

	c.overrideFromOpenSergoConfig(source)
	return nil
}

func (c *OpenSergoConfig) overrideFromOpenSergoConfig(source *OpenSergoConfig) error {
	c.ConfPath = source.ConfPath
	c.Port = source.Port
	return nil
}

const (
	EnvKeyEndpointPort = "OPENSERGO_ENDPOINT_PORT"
)

func (c *OpenSergoConfig) overrideFromSystemEnv() error {
	if portEnv := os.Getenv(EnvKeyEndpointPort); !util.IsBlank(portEnv) {
		port, err := strconv.ParseUint(portEnv, 10, 32)
		if err != nil {
			return err
		} else {
			c.Port = uint(port)
		}
	}
	return nil
}
