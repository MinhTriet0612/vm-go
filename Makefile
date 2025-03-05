
run: main.go 
	go run main.go

hello: 
	echo "Hello, World!"

build: 
	go build -o bin/main main.go


genisoimage:
	genisoimage


all: hello run 
