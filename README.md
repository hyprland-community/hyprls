# HyprLS

<table>
<tr>
	<td> <img src="./demo-completion.png">
	<td> <img src="./demo-hover.png">
</tr>
<tr>
	<td> <img src="./demo-hover-keyword.png">
	<td> <img src="./demo-symbols.png">
	<td> <img src="./demo-colors.png">
</tr>
</table>

A LSP server for Hyprland configuration files.

## Features

Not checked means planned / work in progress.

- [x] Auto-complete
- [x] Hover
  - [ ] TODO: Documentation on hover of categories?
- [x] Go to definition
- [x] Color pickers
- [x] Document symbols
- [ ] Diagnostics
- [ ] Formatting
- [ ] Semantic highlighting

## Installation

### Linux packages

hyprls has packages for various distributions, kindly maintained by other people

[![Packaging status](https://repology.org/badge/vertical-allrepos/hyprls.svg)](https://repology.org/project/hyprls/versions)

### With `go install`

```sh
go install github.com/hyprland-community/hyprls/cmd/hyprls@latest
```

### Pre-built binaries

Binaries for linux are available in [Releases](https://github.com/hyprland-community/hyprls/releases) 

### From source

- Required: [Just](https://just.systems) (`paru -S just` on Arch Linux (btw))

```sh
git clone --recurse-submodules https://github.com/hyprland-community/hyprls
cd hyprls
# installs the binary to ~/.local/bin.
# Make sure that directory exists and is in your PATH
just install
```

## Usage

### With Neovim

_Combine with [The tree-sitter grammar for Hyprlang](https://github.com/tree-sitter-grammars/tree-sitter-hyprlang) for syntax highlighting._

Add this to your `init.lua`:

```lua
-- Hyprlang LSP
vim.api.nvim_create_autocmd({'BufEnter', 'BufWinEnter'}, {
		pattern = {"*.hl", "hypr*.conf"},
		callback = function(event)
				print(string.format("starting hyprls for %s", vim.inspect(event)))
				vim.lsp.start {
						name = "hyprlang",
						cmd = {"hyprls"},
						root_dir = vim.fn.getcwd(),
				}
		end
})
```

### With Emacs
Install [hyprlang-ts-mode](https://github.com/Nathan-Melaku/hyprlang-ts-mode) and [lsp-bridge](https://github.com/manateelazycat/lsp-bridge).

lsp-bridge supports hyprls out of the box, no need to do any configuration.

### VSCode

#### Official Marketplace (VisualStudio Marketplace)

Install it [from the marketplace](https://marketplace.visualstudio.com/items?itemName=gwenn°-lbh.vscode-hyprls).

> [!TIP]
> You can use [the Hyprland extension pack](https://marketplace.visualstudio.com/items?itemName=gwenn°-lbh.hyprland) to also get syntax highlighting.

#### Open VSX (for VSCodium & others)

Install it [on OpenVSX](https://open-vsx.org/extension/gwenn°-lbh/vscode-hyprls)

### Zed

Language server support is provided by the [Hyprlang extension](https://zed.dev/extensions?query=hyprlang).
Detailed installation and setup instructions can be found in the [extension repository](https://github.com/WhySoBad/zed-hyprlang-extension) [maintainer = [@WhySoBad](https://github.com/WhySoBad)].
