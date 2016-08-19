package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	//"strings"
	"strconv"

	"./config"
	"flag"
)

const help string = `Usage:
hydra [OPTIONS] COMMAND

Commands:
  init - Create hydra project [--clean]
  start - Start hydra servers
`

func initialize(c config.Config, clean bool) {
	if clean {
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
			e := exec.Command("cp", "-R", p, base).Run()

			if e != nil {
				log.Fatal(e)
			}

			git := path.Join(base, ".git")
			exec.Command("rm", "-rf", git).Run()

			if service.Config.Dest != "" && service.Config.Src != "" {
				u := path.Join(base, service.Config.Dest)
				cf := service.Config.Src

				err := exec.Command("cp", cf, u).Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func start(conf config.Config) {
	fmt.Println("Starting")
	for i := 0; i < len(conf.Services); i += 1 {
		service := conf.Services[i]
		_, name := path.Split(service.Path)

		env := os.Environ()
		for k, v := range service.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		//cmd := exec.Command(strings.Split(conf.Services[i].Start, " ")...)
		fmt.Println("Starting", name)
		cmd := exec.Command("/usr/bin/env", "node", "./src/server.js", "&")
		cmd.Env = env
		cmd.Dir = path.Join(".hydra", name)
		cmd.Stderr = os.Stderr

		e := cmd.Start()
		fmt.Printf("%d\n", cmd.Process.Pid)
		pid := []byte(strconv.Itoa(cmd.Process.Pid))
		ioutil.WriteFile("pid", pid, 0440)

		if e != nil {
			os.Remove("pid")
			log.Fatal("Error with", name, e)
		}

		fmt.Printf("Started: %s\n", name)
		/*
			e = cmd.Wait()
			if e != nil {
				log.Fatal(e)
			}
		*/
	}
}

func kill(conf config.Config) {
	for i := 0; i < len(conf.Services); i += 1 {
		service := conf.Services[i]
		_, name := path.Split(service.Path)
		pid_file := path.Join(".hype", name, "pid")
		if _, err := os.Stat(pid_file); os.IsExist(err) {
			bytes, e := ioutil.ReadFile(pid_file)
			if e != nil {
				log.Fatal(e)
			}
			pid := string(bytes)
			e = exec.Command("kill", pid).Run()
			if e != nil {
				log.Fatal(e)
			}

			fmt.Printf("Killed %s", name)
		}
	}
}

func main() {
	var c config.Config = config.ReadConfig()
	var clean = flag.Bool("clean", false, "If we want a clean run")
	flag.Parse()

	if cmd := flag.Arg(0); cmd != "" {
		switch cmd {
		case "init":
			initialize(c, *clean)
		case "start":
			start(c)
		case "kill":
			kill(c)
		default:
			fmt.Println(help)
		}
	} else {
		fmt.Println(help)
	}
}
