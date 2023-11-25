lint:
	go fmt ./...
	golangci-lint run ./...

run: build
	./bin/griffon

dot:
	dot -Tpng graph.dot -o output.png

generate:
	# https://github.com/uber-go/mock
	go generate ./...

# command to build this project
build:
	go build -o bin/ ./...