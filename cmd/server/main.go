package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cel-validator/src/lsp"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "stdio", "Communication mode (stdio, tcp)")
	flag.Parse()

	ctx := context.Background()
	server, err := lsp.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	switch mode {
	case "stdio":
		fmt.Fprintln(os.Stderr, "CEL Validator Language Server started in stdio mode")
		server.StartStdio(ctx)
	default:
		log.Fatalf("Unsupported mode: %s", mode)
	}
}
