package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type PasswordAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SSHAuth struct {
	Username       string `yaml:"username"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
}

type SSHConfig struct {
	PasswordAuth `yaml:"passwordAuth"`
	SSHAuth      `yaml:"sshAuth"`
	Ips          []string `yaml:"ips"`
}

var Config []*SSHConfig

func init() {
	// set the config object so that it can be used everywhere and also set the viper config
	// initialize the viper settings
	viper.SetConfigFile(os.Getenv("CONFIG_FILE_PATH"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var ConfigData struct {
		ConfigList []struct {
			SshConfig *SSHConfig `yaml:"sshConfig"`
		} `yaml:"configList"`
	}
	err = viper.UnmarshalExact(&ConfigData)
	if err != nil {
		log.Fatal(err)
	}
	configList := []*SSHConfig{}
	for _, t := range ConfigData.ConfigList {
		configList = append(configList, t.SshConfig)
	}
	Config = configList
}
