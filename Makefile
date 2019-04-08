controller:*.go
	go build -tags gtk_3_12 -o controller
controller.exe:*.go
	go build -tags gtk_3_12 -o controller.exe -ldflags "-H windowsgui"