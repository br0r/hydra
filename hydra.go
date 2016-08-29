package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"flag"
	"github.com/br0r/hydra/config"
)

const BASE_DIR string = ".hydra"
const help string = `Usage:
hydra [OPTIONS] COMMAND

Commands:
  init - Create hydra project [--clean]
  start - Start hydra servers
  stop - Stops started servers
  ls - Show started servers
  logs [name] - Show logs for servers, or for specific server given by name
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
		name := service.Name
		base := path.Join(".hydra/", name)
		if _, err := os.Stat(base); os.IsNotExist(err) {
			fmt.Println("Creating", name)

			var e error
			if strings.Index(p, "git@") == 0 {
				e = exec.Command("git", "clone", p, base).Run()
			} else {
				e = exec.Command("cp", "-R", p, base).Run()
			}

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

			if service.Install != "" {
				args := strings.Split(service.Install, " ")
				cmd := exec.Command(args[0], args[1:]...)

				cmd.Dir = base
				e := cmd.Run()
				if e != nil {
					log.Fatal(e)
				}
			}
		}
	}
}

func start(conf config.Config) {
	for i := 0; i < len(conf.Services); i += 1 {
		service := conf.Services[i]
		name := service.Name

		cmdRunDir := path.Join(".hydra", name)
		pid_file := path.Join(cmdRunDir, "pid")
		log_file_path := path.Join(cmdRunDir, "log")

		if _, err := os.Stat(pid_file); !os.IsNotExist(err) {
			fmt.Printf("%s is already started, run hydra stop if you want to restart\n", name)
			continue
		}

		env := os.Environ()
		for k, v := range service.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Println("Starting", name)
		args := strings.Split(conf.Services[i].Start, " ")
		args = append(args, "2>&1 1> log")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = env
		cmd.Dir = cmdRunDir

		log_file, e := os.Create(log_file_path)
		if e != nil {
			log.Fatal(e)
		}

		_, e = log_file.WriteString(fmt.Sprintf("%s:\n", name))
		if e != nil {
			log.Fatal(e)
		}

		defer log_file.Close()
		cmd.Stdout = log_file
		cmd.Stderr = log_file

		e = cmd.Start()
		pid := []byte(strconv.Itoa(cmd.Process.Pid))
		ioutil.WriteFile(pid_file, pid, 0440)

		if e != nil {
			os.Remove("pid")
			log.Fatal("Error with", name, e)
		}

	}
}

func kill(conf config.Config) {
	for i := 0; i < len(conf.Services); i += 1 {
		service := conf.Services[i]
		name := service.Name
		pid_file := path.Join(".hydra", name, "pid")
		if _, err := os.Stat(pid_file); !os.IsNotExist(err) {
			bytes, e := ioutil.ReadFile(pid_file)
			if e != nil {
				log.Fatalf("Error when reading file %v", e)
			}
			pid := string(bytes)
			exec.Command("kill", pid).Run()
			if e == nil {
				fmt.Printf("Killed %s\n", name)
			}
			os.Remove(pid_file)
		}
	}
}

func ls(conf config.Config) {
	for _, service := range conf.Services {
		name := service.Name
		pid_path := path.Join(BASE_DIR, name, "pid")
		if _, err := os.Stat(pid_path); !os.IsNotExist(err) {
			bytes, e := ioutil.ReadFile(pid_path)
			if e != nil {
				log.Fatal(e)
			}
			pid := string(bytes)
			fmt.Printf("%s running on pid: %s\n", name, pid)
		}
	}
}

func logs(conf config.Config, servers []string) {
	for _, service := range conf.Services {
		if len(servers) > 0 {
			err := false
			for _, name := range servers {
				if name != service.Name {
					err = true
					break
				}
			}
			if err == true {
				continue
			}
		}

		logfile := path.Join(".hydra/", service.Name, "log")
		data, err := ioutil.ReadFile(logfile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
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
		case "stop":
			kill(c)
		case "ls":
			ls(c)
		case "logs":
			logs(c, flag.Args()[1:])
		default:
			fmt.Println(help)
		}
	} else {
		fmt.Println(help)
	}
}
