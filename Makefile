build:
	go build -o bin/go-miner cmd/miner.go

build-win64:
	GOOS=windows GOARCH=amd64 go build -o bin/go-miner-amd64.exe cmd/miner.go
