run: build
	@ ./bin/interpreter

build:
	@ go build -o ./bin/interpreter ./src
