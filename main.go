package main

import (
	"github.com/namsral/flag"
)

type Config struct {
	data string
	bind string
}

func main() {
	var config Config

	flag.StringVar(&config.data, "data", "./", "path to data")
	flag.StringVar(&config.bind, "bind", "0.0.0.0:8000", "[addr]:<port> to bind to")
	flag.Parse()

	ServerInit(config)
}
