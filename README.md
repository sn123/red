# Red 
<img src="https://github.com/sn123/red/raw/master/screenshots/red.png" width="100">

Red is a lightweight webserver for quick file sharing between systems on same network. Red is also a static web server and can help in quickly hosting any folder and making it accessible on the intranet.
It works well for accessing & downloading files from the system to the phone on the same network.

Binaries can be downloaded from ![releases](https://github.com/sn123/red/releases) page.
## How it works
~~Copy over the red server into any folder that should be red-ified.~~
Starting 1.2, Red supports "Red-ifying" any folder by passing it as an argument (see running from source section below).
There's a basic Explorer context menu integration (for now ![Windows only](https://github.com/sn123/red/issues/2)), right-clicking on any folder shows an option to "Redify" which would spawn red serving that folder.
A powershell script is provided as part as 1.2 to make respective registry entries, open powershell:
```powershell
PS c:\red-path> Unblock-File -Path install.ps1
PS c:\red-path> .\install.ps1
```

Run red, once Red starts it will automatically find the next available port and outbound IP address and will print the QR code on console which can be scanned by the phones to navigate to the server.

Since red relies on Golang's built-in behavior, having an index.html file in the root of the folder would cause it to serve the index file instead of showing directory content (https://golang.org/src/net/http/fs.go). It tries to hijack the default html generated to make it more responsive.

![Run](https://github.com/sn123/red/raw/master/screenshots/screenshot.png)

On the phone, point your camera to the QR code

![Phone](https://github.com/sn123/red/raw/master/screenshots/phone.gif)

### Running from source
```bash
$ go build
$ ./red --path=folder-to-redify
```

