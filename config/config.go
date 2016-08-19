package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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

func ReadConfig() Config {
	d, e := ioutil.ReadFile(config_file_name)
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}

	c := Config{}
	err := yaml.Unmarshal(d, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}

	return c
}
