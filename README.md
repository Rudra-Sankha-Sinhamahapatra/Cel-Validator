# CEL Validator for VS Code

A VS Code extension that provides syntax highlighting, validation, and error detection for CEL (Common Expression Language) expressions.

# Demo

https://github.com/user-attachments/assets/bd79a1fd-a0ec-441e-b154-957873f47e0e

## Features

- Syntax highlighting for CEL expressions
- Real-time validation of CEL expressions
- Error diagnostics with detailed messages
- Support for CEL standard library functions
- Works with `.cel` files and embedded CEL in YAML/JSON

## What is CEL?

[Common Expression Language (CEL)](https://github.com/google/cel-spec) is an open source expression language created by Google that's used for policy definition, validation rules, and resource templates. It's used in Kubernetes, Istio, and other cloud-native projects.

## Installation

### Prerequisites

- VS Code 1.75.0 or later
- Go 1.24.1 or later (for building the server)
- Node.js 22.13.1 or later
- npm 10.9.2 or later

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/Rudra-Sankha-Sinhamahapatra/Cel-Validator
   cd Cel-Validator
   ```

2. Build the validation server:
   ```
   go build -o bin/cel-validator cmd/server/main.go
   chmod +x bin/cel-validator
   ```

3. Build the VS Code extension:
   ```
   cd extension
   npm install
   npm run compile
   ```

4. Launch the extension in development mode:
   - Open the project in VS Code
   - Press F5 (or Fn+F5 on Mac)

## Usage

1. Create a file with `.cel` extension (e.g., `expressions.cel`)
2. Write CEL expressions in the file
3. The extension will automatically validate your expressions
4. Invalid expressions will be underlined with error messages

### Example CEL Expressions

```cel
# Valid expressions
5 + 3 * 2
size(["foo", "bar"]) == 2
{"key": "value"}.key == "value"

# Invalid expressions (will show errors)
5 + "string"  
size(123)
```

## Architecture

The CEL Validator consists of two main components:

1. **Go-based Validation Server**:
   - Uses the official Google CEL Go library
   - Provides expression parsing and type checking
   - Reports detailed diagnostic information

2. **VS Code Extension**:
   - Provides syntax highlighting
   - Manages communication with the validation server
   - Displays diagnostics in the editor

## Future Enhancements

- Auto-completion for CEL functions and operators
- Schema-aware validation for YAML and JSON files
- Better positioning of error messages
- Support for multi-line expressions
- Integration with Kubernetes CRDs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Google CEL team for the CEL specification and Go implementation
- VS Code extension development community 
