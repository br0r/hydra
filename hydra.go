package main

import (
  "fmt"
  "os"
  "os/exec"
  "log"
  "io/ioutil"

  "gopkg.in/yaml.v2"
)

const file_name string = ".hydra.yml"

type Config struct {
  Services []struct {
    Git string
    Config struct {
      Src string
      Dest string
    }
  }
}

func main() {
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

  exec.Command("mkdir", ".hydra").Output()

  for i := 0; i < len(c.Services); i += 1 {
    base := ".hydra/" + c.Services[i].Git
    exec.Command("mkdir", base).Output()
    u := base + "/" + c.Services[i].Config.Dest
    fmt.Printf("%s\n", u)
    _, err := exec.Command("touch", u).Output()
    if err != nil {
      log.Fatal(err)
    }
  }
}
