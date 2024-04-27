# Contributing

Thanks for your interest in contributing to HyprLS :3!

## Setup your environment development

You just need to install:

- [Go](https://golang.org/doc/install)
- [Just](https://just.systems)

Then:

```sh
git clone --recurse-submodules https://github.com/your/fork
cd fork
just build
```

Then, you can build a binary locally with `just build`.
To create a "debug build" that logs all the requests, responses and server logs to files (useful for debugging), you can run `just build-debug`.

The debug binary is named `hyprls-debug` and the regular binary is named `hyprls`.

### VSCode

To develop the vscode extension, you'll also need:

- [VSCode](https://code.visualstudio.com/), really useful for debugging the extension, you just have to launch from the "start/debug" menu in the sidebar
- [Bun](https://bun.sh)

Then:

```sh
cd vscode
bun i
```

> **Note:** "Reloading" does not re-build the Go server, you'll need to kill the "Develoment Host Extension" window and kill the terminal that ran the `go:compile` task before launching again.

The VSCode dev environment is set up to use the debug binary in development. The path is absolute and hardcoded, so you may need to change it to where your debug binary is (check `vscode/src/extension.ts`).

## Commit names

We use the [gitmoji](https://gitmoji.dev/) convention for commit names.
