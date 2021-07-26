package main

import (
	"fmt"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func main() {
	fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)

}
