default: run

run:
	# go run main.go add 
	go run main.go add -git https://github.com/zong-zhe/kcl1

build: 
	go build -o build/kpm main.go

test: 
	go test -v ./...	

lint:
	golint ./...

codegen:
	protoc --go_out=. ./kcl.mod.proto
	protoc --gotag_out=auto="toml":. ./kcl.mod.proto
