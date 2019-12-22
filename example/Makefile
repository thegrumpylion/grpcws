
serve: client.wasm wasm_exec.js srv
	@./srv

srv: $(wildcard server/*.go)
	@go build -o srv ./server/

client.wasm: client/main.go
	@GOOS=js GOARCH=wasm go build -o client.wasm ./client/

service/service.pb.go: service/service.proto
	@protoc --go_out=plugins=grpc:. $<

wasm_exec.js:
	@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js .

clean:
	@rm -f client.wasm wasm_exec.js srv