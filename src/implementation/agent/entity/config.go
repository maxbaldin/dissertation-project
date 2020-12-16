package entity

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MaxCores int `yaml:"max_cores"`
	Network  struct {
		PCap struct {
			SnapshotLength int    `yaml:"snapshot_length"`
			BFFFilter      string `yaml:"bff_filter"`
		}
	}
	Integration struct {
		HTTP struct {
			DefaultTimeoutMs int `yaml:"default_timeout_ms"`
		}
		Process struct {
			Repository struct {
				UpdateIntervalMs int `yaml:"update_interval_ms"`
			}
		}
		Collector struct {
			Producer struct {
				URL         string
				QueueLength int `yaml:"queue_length"`
			}
			Aggregator struct {
				InitialBufferLength int `yaml:"initial_buffer_length"`
				FlushIntervalSec    int `yaml:"flush_interval_sec"`
			}
			KnownNodes struct {
				URL               string
				UpdateIntervalSec int      `yaml:"update_interval_sec"`
				Additional        []string `yaml:"additional"`
			} `yaml:"known_nodes"`
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
