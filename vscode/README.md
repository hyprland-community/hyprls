# VSCode Extension for Hyprland configuraiton files

## Installation

Requires [installing `hyprls`](https://github.com/ewen-lbh/hyprls) and having in on your PATH.

## Ignore file preference

By default HyprLS will prefer a `.hyprlsignore` file in your workspace to list filenames or patterns that the language server should ignore. This behavior can be overridden by a configuration option exposed by the VSCode extension.

- **Config key:** `hyprls.preferIgnoreFile` (boolean)
- **Default:** `true` — prefer the `.hyprlsignore` file over the extension setting `hyprls.ignore`.

How to set it in VSCode:

```json
// In your Workspace or User settings (settings.json)
{
	"hyprls.preferIgnoreFile": false,
	"hyprls.ignore": [
		"hyprlock.conf",
		"hypridle.conf"
	]
}
```
