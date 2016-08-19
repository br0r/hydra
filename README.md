# Hydra
CLI tool for starting tests servers for microservices. Used for integration tests.

## Installation

```bash
go get -u github.com/br0r/hydra
```

## Structure for .yml file

*.hydra.yml*

```yml
---

services:
  - name: string #Identifier name.
    path: string #Local path or git repo (only ssh for now).
    start: string #Command to run on hydra start.
    install: string #If we need to do any installation after fetching, specify command here.
    config:
      src: string #Where to find config file to user for service>
      dest: string #Where to put it.
    env: map #If we need any ENV variables set when running start, specify them here
      #E.g NODE_ENV: "test"

```
