package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	qrcode "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/phayes/freeport"
)

func main() {
	port, path := parseFlags()
	if path == "" {
		path = "./"
	}
	log.Printf("Redifying %s", path)
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", hijack(fs))
	ip, port := getOutboundIP(port)
	log.Printf("Listening on %s:%d", ip, port)
	printQrCode(fmt.Sprintf("http://%s:%d", ip, port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func hijack(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := `<html>
		<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			body { margin:0;padding:0; color:#fff;font-size:1.2em;background:#264653}
			body * { box-sizing: border-box;}
			a { display: block; padding: 10px;color:#fff;border-bottom:1px solid #333;margin:0 0 3px;}
			.text-muted {text-align:center;font-size:0.8em;}
			pre { white-space: unset;}
		</style>
		</head><body>
		`
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, r) // call original
		header := recorder.Result().Header.Get("Content-Type")
		w.Header().Set("Content-Type", header)
		if !strings.HasPrefix(header, "text/html") {
			io.Copy(w, recorder.Result().Body)
			return
		}
		// time to hijack the response else we let it go through
		fmt.Fprintf(w, start)
		rsp, _ := ioutil.ReadAll(recorder.Result().Body)
		fmt.Fprintf(w, string(rsp))
		fmt.Fprint(w, "<p class=\"text-muted\">Powered by Red</p></body></html>")
		// not sure why doing a defer close causes writer to go blank?
		recorder.Result().Body.Close()
	})
}

func parseFlags() (int, string) {
	port := flag.Int("port", 0, "port to listen on")
	path := flag.String("path", "", "path to redify")
	flag.Parse()
	return *port, *path
}

func getPort(port int) int {

	// no port passed find a free one and bind
	if port != 0 {
		return port
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

func getOutboundIP(port int) (string, int) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), getPort(port)
}
