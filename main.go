package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	qrcode "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/phayes/freeport"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	ip, port := getOutboundIP()
	log.Printf("Listening on %s:%d", ip, port)
	printQrCode(fmt.Sprintf("http://%s:%d", ip, port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func getPort() int {
	port := flag.Int("port", 0, "port to listen on")
	flag.Parse()
	// no port passed find a free one and bind
	if *port != 0 {
		return *port
	}
	freePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return freePort
}

func printQrCode(content string) {
	obj := qrcode.New2(qrcode.ConsoleColors.BrightBlack, qrcode.ConsoleColors.BrightWhite, qrcode.QRCodeRecoveryLevels.Low)
	obj.Get(content).Print()
}

func getOutboundIP() (string, int) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), getPort()
}
