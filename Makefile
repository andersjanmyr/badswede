
ALL: build-cli

bin:
	mkdir -p bin

build-cli: bin
	go build -o bin/badswede cmds/cli/main.go

run: build-cli
	bin/badswede
