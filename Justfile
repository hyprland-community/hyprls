serverLogsFilepath := `realpath ./logs/server.log || echo ./logs/server.log`
latestTag := `git describe --tags --abbrev=0 || echo commit:$(git rev-parse --short HEAD)`

release tag:
	jq '.version = "{{ tag }}"' < vscode/package.json | sponge vscode/package.json
	sed -i "s/$(grep 'version' default.nix)/  version = \"{{ tag }}\";/" default.nix
	git add vscode/package.json default.nix
	git commit -m "🏷️ Release {{ tag }}"
	git tag -- v{{ tag }}
	cd vscode; bun vsce package; bun vsce publish
	git push
	git push --tags

run:
	just build
	./hyprlang-lsp

build:
	mkdir -p parser/data/sources
	cp hyprland-wiki/pages/Configuring/*.md parser/data/sources/
	go mod tidy
	go build -ldflags "-X main.Version={{ latestTag }}" -o hyprls cmd/hyprls/main.go

build-debug:
	mkdir -p parser/data/sources
	cp hyprland-wiki/pages/Configuring/*.md parser/data/sources/
	go mod tidy
	go build -ldflags "-X main.OutputServerLogs={{ serverLogsFilepath }}" -o hyprlang-lsp cmd/hyprls/main.go

install:
	just build
	mkdir -p ~/.local/bin
	cp hyprls ~/.local/bin/hyprls

pull-wiki:
	git submodule update --init --recursive --remote

parser-data:
	#!/bin/bash
	set -euxo pipefail
	just build
	cd parser/data/generate
	go build -o generator main.go 
	./generator > ../../highlevel.go ast.json
	gofmt -s -w ../../highlevel.go
	jq . < ast.json | sponge ast.json

update-nix-inputs:
	nix flake update
