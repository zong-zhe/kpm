default: run

build: 
	go build -o build/kpm main.go

test: 
	go test -v ./...	

lint:
	golint ./...

codegen:
	protoc --go_out=. ./kcl.mod.proto
	protoc --gotag_out=auto="toml":. ./kcl.mod.proto
