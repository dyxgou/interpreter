# Dyxgou's Programming Language

This programming language is a bit of a wild ride - it's got some *interesting* quirks that make it unsuitable for serious projects. While the code is surprisingly readable, there's zero documentation and some surprising security gaps.

For example:
- Spiral into infinite loops if input parsing isn't *just right*
- Crash spectacularly with stack overflows
- Hide other mysterious behaviors that keep you guessing!

You can see examples of code in the `example` folder.

But here's the thing: despite all its rough edges, working with it has been an absolute blast. There's something genuinely fun about:
- Reverse-engineering how it works from the examples
- Discovering its eccentric personality through trial and error
- The thrill when something actually works as expected!

*Inspired by the fantastic book [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball*

## How to install

1) Download the repository:
```sh
$ git clone https://github.com/dyxgou/interpreter
```

2) Compile the project with `make`.
This will open an REPL for you to insert any command you like
```sh
$ make
```

Compile the project and execute any file you want
```sh
$ make execute FILE=/path/to/file
```
