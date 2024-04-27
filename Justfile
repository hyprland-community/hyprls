serverLogsFilepath := `realpath ./logs/server.log`

run:
	just build
	./hyprlang-lsp

build:
	mkdir -p parser/data/sources
	cp hyprland-wiki/pages/Configuring/*.md parser/data/sources/
	go mod tidy
	go build -ldflags "-X main.OutputServerLogs={{ serverLogsFilepath }}" -o hyprlang-lsp cmd/main.go

install:
	just build
	cp hyprlang-lsp ~/.local/bin/hyprls

parser-data:
	#!/bin/bash
	set -euxo pipefail
	just build
	cd parser/data/generate
	go build -o generator main.go 
	./generator > ../../highlevel.go ast.json
	gofmt -s -w ../../highlevel.go
	jq . < ast.json | sponge ast.json
