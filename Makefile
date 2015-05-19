
ALL: build-cli build-web

bin:
	mkdir -p bin

.PHONY: build-cli
build-cli: bin
	go build -o bin/badswede cmds/cli/main.go

.PHONY: run
run: build-cli
	bin/badswede

.PHONY: build-web
build-web: bin
	go build -o bin/badswede-web cmds/web/main.go

.PHONY: web
web: build-web
	(cd cmds/web/ && ../../bin/badswede-web)

.PHONY: clean
clean:
	rm -rf bin
