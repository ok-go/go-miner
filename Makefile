build:
	go build -o bin/go-miner go-miner/cmd

build-win64:
	GOOS=windows GOARCH=amd64 go build -o bin/go-miner-amd64.exe .
