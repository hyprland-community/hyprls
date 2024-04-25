run:
	just build
	./hyprlang-lsp

build:
	cp hyprland-wiki/pages/Configuring/Variables.md parser/data/
	go mod tidy
	go build -o hyprlang-lsp cmd/main.go

install:
	just build
	cp hyprlang-lsp ~/.local/bin

parser-data:
	#!/bin/bash
	set -euxo pipefail
	cd parser/data/generate
	go build -o generator main.go 
	./generator < ../../../hyprland-wiki/pages/Configuring/Variables.md > ../../highlevel.go ast.json
	jq . < ast.json | sponge ast.json
