.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -a -tags netgo -ldflags "-s -w" -o pdf.wasm
	ls -lah pdf.wasm

.PHONY: copy-js
copy-js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

.PHONY: serve
serve:
	go run server.go
