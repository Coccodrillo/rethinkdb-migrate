package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type Config struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	CertFile string `yaml:"cert_file"`
}

func NewConfig() *Config {
	c := &Config{}
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return c
}
