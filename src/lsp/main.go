package lsp

import (
	"context"
	"fmt"
	"os"

	"github.com/google/cel-go/cel"
	protocol "go.lsp.dev/protocol"
)

// Server represents the CEL validator language server
type Server struct {
	env *cel.Env
}

// NewServer creates a new language server for CEL validation
func NewServer() (*Server, error) {
	// Creating a default CEL environment
	env, err := cel.NewEnv(cel.StdLib())
	if err != nil {
		return nil, err
	}

	return &Server{
		env: env,
	}, nil
}

// StartStdio starts the language server using standard input/output
func (s *Server) StartStdio(ctx context.Context) error {
	// A very simple implementation that just validates any CEL expressions
	// in text files and prints diagnostics to stderr for demonstration

	fmt.Fprintln(os.Stderr, "CEL Validator started in stdio mode")
	fmt.Fprintln(os.Stderr, "Type a CEL expression and press Enter to validate:")

	// Simple read loop
	var input string
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if _, err := fmt.Scanln(&input); err != nil {
				if err.Error() == "EOF" {
					return nil
				}
				continue
			}

			diagnostics := s.validateCEL(input)

			// Print diagnostics
			if len(diagnostics) > 0 {
				for _, diag := range diagnostics {
					fmt.Fprintf(os.Stderr, "Error: %s\n", diag.Message)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Valid CEL expression!")
			}
		}
	}
}

// validateCEL validates a CEL expression
func (s *Server) validateCEL(content string) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	ast, issues := s.env.Parse(content)

	if issues != nil && issues.Err() != nil {
		// Diagnostic for the syntax error
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: uint32(len(content))},
			},
			Severity: protocol.DiagnosticSeverityError,
			Source:   "cel-validator",
			Message:  issues.Err().Error(),
		})
	} else {
		_, issues = s.env.Check(ast)
		if issues != nil && issues.Err() != nil {
			// Diagnostic for the type error
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: uint32(len(content))},
				},
				Severity: protocol.DiagnosticSeverityError,
				Source:   "cel-validator",
				Message:  issues.Err().Error(),
			})
		}
	}

	return diagnostics
}
