run:
	just build
	./hyprlang-lsp

build:
	go mod tidy
	go build -o hyprlang-lsp .

install:
	just build
	cp hyprlang-lsp ~/.local/bin
