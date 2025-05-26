serverLogsFilepath := `realpath ./logs/server.log || echo ./logs/server.log`
latestTag := `git describe --tags --abbrev=0 || echo commit:$(git rev-parse --short HEAD)`
latestVersion := `git describe --tags --abbrev=0  | sed 's/v//' || echo commit:$(git rev-parse --short HEAD)`
# Parse content of https://raw.githubusercontent.com/hyprwm/hyprland-wiki/main/pages/version-selector.md to get latest documented version

check-for-hyprland-updates:
	#!/bin/env bash
	set -euxo pipefail
	latestHyprlandVersion=`curl -s https://raw.githubusercontent.com/hyprwm/hyprland-wiki/main/pages/version-selector.md | grep -oP 'v\d+\.\d+\.\d+' | head -n 1 | sed 's/v//'`
	touch hyprland_version
	if [ "$(cat hyprland_version)" != "$latestHyprlandVersion" ]; then
		echo New version $latestHyprlandVersion released!!! update time :3
		echo "$latestHyprlandVersion" > hyprland_version
		just pull-wiki
		just parser-data
	else
		echo Nyathing to update
	fi

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
	mkdir -p logs
	touch logs/server.log
	mkdir -p parser/data/sources
	# cp hyprland-wiki/pages/Configuring/*.md parser/data/sources/
	go mod tidy
	go build -ldflags "-X main.HyprlandWikiVersion=$(cat hyprland_version) -X main.HyprlsVersion={{ latestVersion }}" -o hyprls cmd/hyprls/main.go

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
	#!/bin/bash
	git submodule update --init --recursive --remote
	cd hyprland-wiki 
	hash=$(git log --all --oneline --grep="versions: add $(cat ../hyprland_version)" | cut -d' ' -f1)
	echo Using wiki https://github.com/hyprwm/hyprland-wiki/commit/$hash
	git checkout $hash
	cd ..
	cp hyprland-wiki/pages/Configuring/*.md parser/data/sources/

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
