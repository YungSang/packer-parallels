package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/yungsang/packer-parallels/builder/parallels/iso"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(iso.Builder))
	server.Serve()
}
