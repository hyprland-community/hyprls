/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import { commands, ExtensionContext, workspace } from "vscode"

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node"

let client: LanguageClient

export function activate(context: ExtensionContext) {
  const serverModule = "hyprls"

  // If the extension is launched in debug mode then the debug server options are used
  // Otherwise the run options are used
  const serverOptions: ServerOptions = {
    run: {
      command: serverModule,
      transport: TransportKind.stdio,
    },
    debug: {
      command: "/home/uwun/projects/hyprls/hyprlang-lsp",
      transport: TransportKind.stdio,
    },
  }

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for plain text documents
    documentSelector: [{ scheme: "file", language: "hyprlang" }],
    outputChannelName: "HyprLS",
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("*.hl"),
      configurationSection: "hyprls",
    },
  }

  context.subscriptions.push(
    commands.registerCommand("vscode-hyprls.restart-lsp", () => {
      client.restart()
    })
  )

  // Create the language client and start the client.
  client = new LanguageClient("hyprlang", "Hypr", serverOptions, clientOptions)

  // Start the client. This will also launch the server
  client.start()
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined
  }
  return client.stop()
}
