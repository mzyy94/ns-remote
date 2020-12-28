package main

import (
	"flag"

	"ns-remote/server"

	"github.com/mzyy94/nscon"
)

func main() {
	var (
		device = flag.String("device", "/dev/hidg0", "simulating hid gadget path")
	)
	controller := nscon.NewController(*device)

	server.StartHTTPServer(controller)
}
