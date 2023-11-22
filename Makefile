lint:
	go fmt ./...
	golangci-lint run ./...

run:
	go run config.go \
		config_body_schema.go \
		config_spec.go \
		main.go \
		parser.go \
		parser_body_schema.go \
		parser_spec.go \
		utils.go \
		graph.go 

dot:
	dot -Tpng graph.dot -o output.png
