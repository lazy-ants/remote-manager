package main

import (
	"github.com/lazy-ants/remote-manager/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
