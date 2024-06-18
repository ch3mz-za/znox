build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o znox.exe

run:
	go run main.go