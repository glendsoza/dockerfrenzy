package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type PasswordAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SSHAuth struct {
	Username       string `yaml:"username"`
	PrivateKeyFile string `yaml:"privateKeyFile"`
}

type SSHConfig struct {
	PasswordAuth `yaml:"passwordAuth"`
	SSHAuth      `yaml:"sshAuth"`
	Ips          []string `yaml:"ips"`
}

type Config struct {
	config []*SSHConfig
	lock   sync.Mutex
}

func NewConfig() (*Config, error) {
	config := &Config{config: []*SSHConfig{}}
	err := config.load()
	return config, err
}

func (c *Config) Get() []*SSHConfig {
	// create a copy and return it
	c.lock.Lock()
	defer c.lock.Unlock()
	configList := []*SSHConfig{}
	for _, c := range c.config {
		temp := *c
		configList = append(configList, &temp)
	}
	return configList
}

func (c *Config) load() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	var ConfigData struct {
		ConfigList []struct {
			SshConfig *SSHConfig `yaml:"sshConfig"`
		} `yaml:"configList"`
	}
	err = viper.UnmarshalExact(&ConfigData)
	if err != nil {
		return err
	}
	configList := []*SSHConfig{}
	for _, t := range ConfigData.ConfigList {
		configList = append(configList, t.SshConfig)
	}
	c.config = configList
	return nil
}

func (c *Config) Reload() error {
	return c.load()
}

func (c *Config) Update(data []byte) error {
	// write the data to the file
	if err := os.WriteFile(fmt.Sprintf("%s/config.yaml", os.Getenv("CONFIG_FOLDER")), data, 0600); err != nil {
		return err
	}
	return c.load()
}

func init() {
	// set the config object so that it can be used everywhere and also set the viper config
	// initialize the viper settings
	// if the config file is not present create a empty yaml file
	configFile := fmt.Sprintf("%s/config.yaml", os.Getenv("CONFIG_FOLDER"))
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(fmt.Sprintf("%s/config.yaml", os.Getenv("CONFIG_FOLDER")), []byte(""), 0600)

	}
	viper.SetConfigFile(configFile)
}
