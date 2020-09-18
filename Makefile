all: linux windows

linux:
	go build -o bin/server -ldflags="-s -w" server/server.go
	go build -o bin/client -ldflags="-s -w" client/client.go
	go build -o bin/tm-server -ldflags="-s -w" tm-server/server.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe -ldflags="-s -w" server/server.go
	GOOS=windows GOARCH=amd64 go build -o bin/client.exe -ldflags="-s -w" client/client.go
	GOOS=windows GOARCH=amd64 go build -o bin/tm-server.exe -ldflags="-s -w" tm-server/server.go

.PHONY: all linux windows
