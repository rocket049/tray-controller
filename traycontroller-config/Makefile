traycontroller-config:*.go
	go build -o traycontroller-config -ldflags -s
	cp traycontroller-config ../deb/usr/local/bin/
traycontroller-config.exe:*.go
	go build -ldflags "-s -H windowsgui"