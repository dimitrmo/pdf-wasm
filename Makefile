.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o pdf.wasm
	ls -lah pdf.wasm

.PHONY: copy-js
copy-js:
	rm -f wasm_exec.js
	wget https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js

.PHONY: serve
serve:
	go run server.go

.PHONY: run
run: build copy-js serve
