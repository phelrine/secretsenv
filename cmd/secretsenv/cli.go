package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/phelrine/secretsenv"
	"gopkg.in/yaml.v3"
)

type CLI struct {
	ConfigPath string
	Args       []string
	Config     Config
	Loaders    map[string]secretsenv.SecretLoader
}

type Config map[string]ConfigItem

// ConfigItem represents a configuration item
//
// It contains the type of the loader, the secret ID, the loader options,
// and the mapping of the secret values to the environment variables.
type ConfigItem struct {
	Type     string                  `yaml:"type,omitempty"`
	SecretId string                  `yaml:"secretId"`
	Option   secretsenv.SecretOption `yaml:"option"`
	Mapping  map[string]*string      `yaml:",inline"`
}

func (c *CLI) Parse() error {
	flag.StringVar(&c.ConfigPath, "config", "", "path to config file")
	flag.Parse()
	if c.ConfigPath == "" {
		// if not specified, search for .secretsenv.yml in the parent directories
		configPath, err := c.searchConfig()
		if err != nil {
			return err
		}
		c.ConfigPath = configPath
	}
	c.Args = flag.Args()
	return nil
}

func (c *CLI) Load() error {
	configData, err := os.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configData, &c.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *CLI) Run() error {
	var configs []string
	if len(c.Args) != 0 {
		configs = c.Args
	} else {
		for name := range c.Config {
			configs = append(configs, name)
		}
	}
	for _, arg := range configs {
		item, ok := c.Config[arg]
		if !ok {
			return fmt.Errorf(
				"specified secret %s not found in the configuration file",
				arg,
			)
		}
		loader, ok := c.Loaders[item.Type]
		if !ok {
			return fmt.Errorf(
				"specified loader %s not found",
				item.Type,
			)
		}
		secrets, err := loader.Load(item.SecretId, item.Option)
		if err != nil {
			return err
		}
		// map the secret values to the environment variables
		mappingResult, err := c.mapSecrets(item, secrets)
		if err != nil {
			return err
		}
		for key, value := range mappingResult {
			fmt.Printf("export %s=%s\n", key, value)
		}
	}
	return nil
}

// searchConfig searches for the .secretsenv.yml file in the parent directories
//
// It starts from the current working directory and goes up the directory tree
// until it finds the config file or reaches the root directory.
func (c *CLI) searchConfig() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		configPath := filepath.Join(dir, ".secretsenv.yml")
		_, err := os.Stat(configPath)
		if err == nil {
			return configPath, nil
		}
		newDir := filepath.Dir(dir)
		if newDir == dir {
			return "", fmt.Errorf("config file not found")
		}
		dir = newDir
	}
}

// mapSecrets maps the secret values to the environment variables
//
// It uses the mapping defined in the configuration file to map the secret values
// to the environment variables. If a mapping is not defined, it uses the variable
// name as the key.
func (c *CLI) mapSecrets(item ConfigItem, secrets secretsenv.SecretResult) (secretsenv.SecretMapping, error) {
	mappingResult := make(secretsenv.SecretMapping)
	interpolator := &secretsenv.Interpolator{}
	for variableName, secretKey := range item.Mapping {
		if secretKey == nil {
			value, ok := secrets[variableName]
			if !ok {
				return nil, fmt.Errorf(
					"specified key %s not found in the secret %s",
					variableName,
					item.SecretId,
				)
			}
			mappingResult[variableName] = value
			continue
		}
		value, err := interpolator.Interpolate(*secretKey, secrets)
		if err != nil {
			return nil, err
		}
		mappingResult[variableName] = value
	}
	return mappingResult, nil
}
