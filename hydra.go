package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"flag"

	"gopkg.in/yaml.v2"
)

const help string = `Usage:
hydra COMMAND [OPTIONS]

Commands:
  init - Create hydra project [--clean]
  start - Start hydra servers
`

const file_name string = ".hydra.yml"

type Config struct {
	Services []struct {
		Path   string
		Start  string
		Config struct {
			Src  string
			Dest string
		}
		Env map[string]string
	}
}

func initialize() {
	var clean = flag.Bool("clean", false, "If we want a clean run")
	flag.Parse()

	d, e := ioutil.ReadFile(file_name)
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

	if *clean {
		_, e := exec.Command("rm", "-rf", ".hydra").Output()
		if e != nil {
			log.Fatal(e)
		}
	}

	exec.Command("mkdir", ".hydra").Output()

	for i := 0; i < len(c.Services); i += 1 {
		service := c.Services[i]
		p := service.Path
		_, name := path.Split(p)
		base := path.Join(".hydra/", name)
		if _, err := os.Stat(base); os.IsNotExist(err) {
			fmt.Println("Creating", name)
			_, e := exec.Command("cp", "-R", p, base).Output()
			if e != nil {
				log.Fatal(e)
			}

			if service.Config.Dest != "" && service.Config.Src != "" {
				u := path.Join(base, service.Config.Dest)
				cf := service.Config.Src

				_, err := exec.Command("cp", cf, u).Output()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func start() {
	fmt.Println("Starting")
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "init":
			initialize()
		case "start":
			start()
		default:
			fmt.Println(help)
		}
	} else {
		fmt.Println(help)
	}
}
