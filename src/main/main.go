package main

import (
	"flag"
	"server"
)

func main() {

	var port int

	flag.IntVar(&port, "Port", 8085, "Port the server listens to")

	flag.Parse()

	server.Run(uint16(port))

}
