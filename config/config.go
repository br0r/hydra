package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const config_file_name string = ".hydra.yml"

type Config struct {
	Services []struct {
		Name   string
		Path   string
		Start  string
		Config struct {
			Src  string
			Dest string
		}
		Env     map[string]string
		Install string
	}
}

func ReadConfig() (Config, error) {
	c := Config{}
	d, e := ioutil.ReadFile(config_file_name)
	if e != nil {
		return c, e
	}

	err := yaml.Unmarshal(d, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
