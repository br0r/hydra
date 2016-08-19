package main

import (
  "fmt"
  "path"
  "os"
  "os/exec"
  "log"
  "io/ioutil"

  "flag"

  "gopkg.in/yaml.v2"
)


const file_name string = ".hydra.yml"

type Config struct {
  Services []struct {
    Path string
    Start string
    Config struct {
      Src string
      Dest string
    }
  }
}

func main() {
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
    p := c.Services[i].Path
    _, name := path.Split(p)
    base := path.Join(".hydra/", name)
    if _, err := os.Stat(base); os.IsNotExist(err) {
      fmt.Println("Creating", name)
      _, e := exec.Command("cp", "-R", p, base).Output()
      if e != nil {
        log.Fatal(e)
      }

      u := path.Join(base, c.Services[i].Config.Dest)
      fmt.Printf("%s\n", u)
      _, err := exec.Command("touch", u).Output()
      if err != nil {
        log.Fatal(err)
      }
    }
  }
}
