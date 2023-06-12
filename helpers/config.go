package helpers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Service struct {
        Port string `yaml:"port"`
    } `yaml:"service"`
    Database struct {
		Name string `yaml:"name"`
		Password string `yaml:"password"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Username string `yaml:"username"`
	} `yaml:"database"`
    Rabbitmq struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"rabbitmq"`
}

func (c *Config) ReadConf(f string) (*Config, error) {
    buf, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(buf, c)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %w", f, err)
    }
    return c, err
}