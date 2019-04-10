controller:*.go
	go build -tags gtk_3_12 -o controller -ldflags -s
controller.exe:*.go
	go build -tags gtk_3_12 -o controller.exe -ldflags "-s -H windowsgui"

