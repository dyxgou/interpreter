run: build
	@ ./bin/interpreter

execute: build_execute
	@ ./bin/executer $(ARGS)


build_execute:
	@ go build -o ./bin/executer ./cmd/executer/execute.go


build:
	@ go build -o ./bin/interpreter ./cmd/interpreter/main.go
