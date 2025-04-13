run: build
	@ ./bin/interpreter

execute: build_execute
	@ ./bin/executer $(FILE)


build_execute:
	@ go build -o ./bin/executer ./cmd/executer/main.go


build:
	@ go build -o ./bin/interpreter ./cmd/interpreter/main.go
