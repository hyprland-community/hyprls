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
						settings = {
							hyprls = {
								preferIgnoreFile = true, -- set to false to prefer `hyprls.ignore`
								ignore = {"hyprlock.conf", "hypridle.conf"}
							}
						}
					}
		end
})
```

You can control whether HyprLS prefers a workspace `.hyprlsignore` file or the editor settings with the `hyprls.preferIgnoreFile` option. Example configurations:

- Using `vim.lsp.start` (example above) — set `settings.hyprls.preferIgnoreFile` to `false` to force the server to use `settings.hyprls.ignore`.

- Using `nvim-lspconfig`:

```lua
local lspconfig = require('lspconfig')
lspconfig.hyprlang.setup{
	cmd = {"hyprls"},
	settings = {
		hyprls = {
			preferIgnoreFile = false,
			ignore = {"hyprlock.conf", "hypridle.conf"}
		}
	}
}
```

When `preferIgnoreFile` is `true` (the default), HyprLS will read `.hyprlsignore` from your workspace root. When it's `false`, it will use the `hyprls.ignore` array from your editor configuration instead.

Example `.hyprlsignore` (create this file at the workspace root):

```
# ignore session-specific files
hyprlock.conf
hypridle.conf
# ignore any file named workspace-specific.ignore
workspace-specific.ignore
```

Notes for Neovim users:

- If you set `preferIgnoreFile = true`, HyprLS will use the workspace `.hyprlsignore` file and ignore any `settings.hyprls.ignore` values passed from Neovim.
- If you set `preferIgnoreFile = false`, HyprLS will use the `ignore` list you provide in `settings.hyprls` (see `nvim-lspconfig` example above).

### With Emacs
Language server support is provided by the [lsp-bridge](https://github.com/manateelazycat/lsp-bridge).

Just install [lsp-bridge](https://github.com/manateelazycat/lsp-bridge) in Emacs, that's all, no need to do any configuration.

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

## Configuration

### Ignoring some files

_Thanks to [@sansmoraxz](https://github.com/sansmoraxz) for this feature ^^_

By default, HyprLS ignores `hyprlock.conf` and `hypridle.conf` files, since they aren't supported yet.

You can create a `.hyprlsignore` file that lists filenames HyprLS should not run on. Files are relative to the workspace root, which is determined by your IDE (for example, for VSCode, it's the folder you opened it with)