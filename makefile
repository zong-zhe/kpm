default: run

add:
	go run main.go add -git https://github.com/zong-zhe/kcl1

init:
	go run main.go init mykcl

build: 
	go build -o /Users/shijun/go/bin/kpm main.go

test: 
	go test -v ./...	

lint:
	golint ./...

codegen:
	protoc --go_out=. ./kcl.mod.proto
	protoc --gotag_out=auto="toml":. ./kcl.mod.proto
