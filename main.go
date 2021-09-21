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

var fallback = false
var rootPath http.Dir

func main() {
	port, path, fallbackResource := parseFlags()
	if path == "" {
		path = "./"
	}
	fallback = fallbackResource
	log.Printf("Redifying %s", path)
	rootPath = http.Dir(path)
	fs := http.FileServer(rootPath)
	http.Handle("/", hijack(fs))
	ip, port := getOutboundIP(port)
	log.Printf("Listening on http://%s:%d", ip, port)
	printQrCode(fmt.Sprintf("http://%s:%d", ip, port))
	log.Printf("Listening on http://localhost:%d", port)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()
	var input string
	fmt.Println("Press enter/return to exit")
	fmt.Scanln(&input)
}

func hijack(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := `
		<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			body { margin:0;padding:0; color:#fff;font-size:1.2em;background:#264653}
			body * { box-sizing: border-box;}
			a { display: block; padding: 10px;color:#fff;border-bottom:1px solid #333;margin:0 0 3px;}
			.text-muted {text-align:center;font-size:0.8em;}
			pre { white-space: unset;}
			p.text-muted a {display: inline-block;padding:5px;}
		</style>
		</head><body>
		`
		end := `
			<p class="text-muted">Powered by <a href="https://github.com/sn123/red">Red</a></p>
			</body>
		`
		f, err := rootPath.Open(r.URL.Path)
		if err == nil {
			// we found the file do nothing
			log.Printf("found the file %s", r.URL.Path)
			f.Close()
		}
		if err != nil && os.IsNotExist(err) && fallback == true {
			log.Print("switching to fallback")
			r.URL.Path = "/"
		}
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, r) // call original
		header := recorder.Result().Header.Get("Content-Type")
		w.Header().Set("Content-Type", header)
		if !strings.HasPrefix(header, "text/html") {
			io.Copy(w, recorder.Result().Body)
			recorder.Result().Body.Close()
			return
		}
		rsp, _ := ioutil.ReadAll(recorder.Result().Body)
		stringResponse := string(rsp)
		isSnippet := strings.HasPrefix(stringResponse, "<pre>")
		if !isSnippet {
			fmt.Fprint(w, string(rsp))
			recorder.Result().Body.Close()
			return
		}

		fmt.Fprint(w, start)
		fmt.Fprint(w, string(rsp))
		fmt.Fprint(w, end)
		// not sure why doing a defer close causes writer to go blank?
		recorder.Result().Body.Close()
	})
}

func parseFlags() (int, string, bool) {
	port := flag.Int("port", 0, "port to listen on")
	path := flag.String("path", "", "path to redify")
	fallbackResource := flag.Bool("fs", false, "fallback to /index.html")
	flag.Parse()
	return *port, *path, *fallbackResource
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

func doNothing() { // Noncompliant
}
