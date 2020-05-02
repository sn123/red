# Red 
<img src="https://github.com/sn123/red/raw/master/screenshots/red.png" width="100">

Red is a lightweight webserver for quick file sharing between systems on same network. Red is also a static web server and can help in quickly hosting any folder and making it accessible on the intranet.
It works well for accessing & downloading files from the system to the phone on the same network.
## How it works
Copy over the red server into any folder that should be red-ified. Run red, once Red starts it will automatically find the next available port and outbound IP address and will print the QR code on console which can be scanned by the phones to navigate to the server.

Since red relies on Golang's built-in behavior, having an index.html file in the root of the folder would cause it to serve the index file instead of showing directory content (https://golang.org/src/net/http/fs.go).

![Run](https://github.com/sn123/red/raw/master/screenshots/screenshot.png)

### Running from source
```bash
$ go build
$ ./red
```
