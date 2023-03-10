run:
	go run main.go

run_init:
	go run main.go init test_name

build: 
	go build -o build/kpm main.go

test: 
	go test -v ./...	

codegen:
	protoc --go_out=gen ./kcl.mod.proto


