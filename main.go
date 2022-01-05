package main

import (
  "log"
  "syscall"
)

func main() {
  const bin = "/perm/grafana/bin/grafana-server"
  if err := syscall.Exec(bin, []string{bin, "-homepath=/perm/grafana"}, nil); err != nil {
    log.Fatal(err)
  }
}
