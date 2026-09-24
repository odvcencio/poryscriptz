const vscode = require('vscode');
const { LanguageClient } = require('vscode-languageclient/node');

let client;

async function activate(context) {
  const binary = vscode.workspace.getConfiguration('poryz').get('server.path', 'poryz');
  const serverOptions = { command: binary, args: ['lsp'] };
  const clientOptions = { documentSelector: [{ scheme: 'file', language: 'poryz' }] };
  client = new LanguageClient('poryz', 'poryz language server', serverOptions, clientOptions);
  context.subscriptions.push(client);
  await client.start();
}

async function deactivate() {
  if (client) await client.stop();
}

module.exports = { activate, deactivate };
