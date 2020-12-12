package entity

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		ListenAddr string `yaml:"listen_addr"`
	}
	KnownNodes struct {
		UpdateIntervalSec int `yaml:"update_interval_sec"`
	} `yaml:"known_nodes"`
	Integration struct {
		Db struct {
			ConnectionString string `yaml:"connection_string"`
		}
	}
}

func (c *Config) FromFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}
