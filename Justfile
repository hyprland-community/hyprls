build:
	go mod tidy
	go build main.go -o hyprlang-lsp

install:
	just build
	cp hyprlang-lsp ~/.local/bin
