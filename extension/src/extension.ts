import * as path from 'path';
import * as vscode from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
    let serverPath = vscode.workspace.getConfiguration('celValidator').get<string>('serverPath');

    if(!serverPath) {
        serverPath = context.asAbsolutePath(
            path.join('server', process.platform === 'win32' ? 'cel-validator.exe' : 'cel-validator')
        );
    }

    const serverOptions: ServerOptions = {
        command: serverPath,
        args:["--mode","stdio"],
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            {scheme:'file',language:'cel'},
            {scheme:'file',language:'yaml'},
            {scheme:'file', language:'json'},
        ],
        synchronize:{
            //notify the server about file changes to '.cel' and '.yaml' files contained in the workspace
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.{cel,yaml,json}'),
        },
    };

    client = new LanguageClient (
        'celValidator',
        'CEL Validator',
        serverOptions,
        clientOptions
    );

    client.start();
    context.subscriptions.push(client)
}

export function deactivate(): Thenable<void> | undefined {
    if(!client) {
        return undefined;
    }
    return client.stop();
}

