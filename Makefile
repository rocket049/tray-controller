controller:*.go
	go build -tags gtk_3_12 -o controller -ldflags -s
controller.exe:*.go
	go build -tags gtk_3_12 -o controller.exe -ldflags "-s -H windowsgui"
	
traycontroller-1.0.0-ubuntu_amd64.deb:controller
	cp controller ./deb/usr/local/bin/traycontroller
	dpkg -b ./deb/ traycontroller-1.0.0-ubuntu_amd64.deb
