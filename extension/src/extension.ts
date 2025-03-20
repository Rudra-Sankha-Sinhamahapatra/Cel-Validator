import * as path from 'path';
import * as vscode from 'vscode';
import * as fs from "fs";
import * as cp from 'child_process';


const diagnosticCollection = vscode.languages.createDiagnosticCollection('cel-validator');

export function activate(context: vscode.ExtensionContext) {
    console.log('CEL Validator extension activated');
    
    let serverPath = vscode.workspace.getConfiguration('celValidator').get<string>('serverPath');
    if(!serverPath) {
        if (fs.existsSync(path.join(context.extensionPath, '..', 'bin', 'cel-validator'))) {
            serverPath = path.join(context.extensionPath, '..', 'bin', 'cel-validator');
        } else {
            serverPath = context.asAbsolutePath(path.join('server', 'cel-validator'));
        }
        console.log("Using server path: " + serverPath);
    }

    // Createing an output channel
    const outputChannel = vscode.window.createOutputChannel("CEL Validator");
    context.subscriptions.push(outputChannel);
    
    // Registering the diagnostics collection
    context.subscriptions.push(diagnosticCollection);
    

    function validateDocument(document: vscode.TextDocument) {
        if (document.languageId !== 'cel') return;
        

        if (!serverPath || !fs.existsSync(serverPath)) {
            outputChannel.appendLine(`Server binary not found at: ${serverPath}`);
            outputChannel.show();
            return;
        }
        
        outputChannel.appendLine(`Validating: ${document.uri.fsPath}`);
        outputChannel.appendLine(`Content: ${document.getText()}`);
        
        try {
            // Using spawn with specific options to avoid TypeScript errors
            const process = cp.spawn(serverPath, ['--mode', 'stdio'], {
                stdio: ['pipe', 'pipe', 'pipe']
            });
            
            let output = '';
            
            // Writing the document content to the server
            if (process.stdin) {
                process.stdin.write(document.getText());
                process.stdin.end();
            }
            

            if (process.stdout) {
                process.stdout.on('data', (data: Buffer) => {
                    const text = data.toString();
                    output += text;
                    outputChannel.appendLine(`Server stdout: ${text}`);
                });
            }
            
            if (process.stderr) {
                process.stderr.on('data', (data: Buffer) => {
                    const text = data.toString();
                    output += text;
                    outputChannel.appendLine(`Server stderr: ${text}`);
                });
            }
            
            // Process exit
            process.on('close', (code: number | null) => {
                outputChannel.appendLine(`Server process exited with code ${code}`);
                
                // Clearing previous diagnostics
                diagnosticCollection.delete(document.uri);
                
                if (output.includes("Error:") || code !== 0) {
                    const diagnostic = new vscode.Diagnostic(
                        new vscode.Range(0, 0, document.lineCount, 0),
                        output || "CEL validation failed",
                        vscode.DiagnosticSeverity.Error
                    );
                    diagnosticCollection.set(document.uri, [diagnostic]);
                } else {
                    outputChannel.appendLine("CEL expression is valid");
                }
            });
            
            process.on('error', (err: Error) => {
                outputChannel.appendLine(`Server error: ${err.message}`);
            
                outputChannel.show();
                
                const diagnostic = new vscode.Diagnostic(
                    new vscode.Range(0, 0, document.lineCount, 0),
                    `Failed to run CEL validator: ${err.message}`,
                    vscode.DiagnosticSeverity.Error
                );
                diagnosticCollection.set(document.uri, [diagnostic]);
            });
        } catch (err) {
            outputChannel.appendLine(`Error running validator: ${err}`);
            outputChannel.show();
        }
    }
    
    //  handlers for document changes
    context.subscriptions.push(
        vscode.workspace.onDidOpenTextDocument(validateDocument),
        vscode.workspace.onDidChangeTextDocument(event => validateDocument(event.document)),
        vscode.workspace.onDidSaveTextDocument(validateDocument)
    );
    
    vscode.workspace.textDocuments.forEach(validateDocument);

    const disposable = vscode.commands.registerCommand('celValidator.showOutput', () => {
        outputChannel.show();
    });
    context.subscriptions.push(disposable);
}

export function deactivate() {
    diagnosticCollection.clear();
    diagnosticCollection.dispose();
}

