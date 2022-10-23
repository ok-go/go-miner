build:
	go build -v -o bin/go-miner go-miner/cmd

build-win:
	go build -v -o bin/go-miner.exe -ldflags -H=windowsgui go-miner/cmd
