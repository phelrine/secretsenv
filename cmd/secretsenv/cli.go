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
	SecretsEnv *secretsenv.SecretsEnv
}

type Config map[string]secretsenv.SecretOption

func NewCLI(loaders map[string]secretsenv.SecretLoader) *CLI {
	return &CLI{
		SecretsEnv: &secretsenv.SecretsEnv{
			Loaders: loaders,
		},
	}
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
		result, err := c.SecretsEnv.Load(item)
		if err != nil {
			return err
		}
		for key, value := range result {
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
