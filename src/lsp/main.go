package lsp

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/cel-go/cel"
	protocol "go.lsp.dev/protocol"
)

// Server represents the CEL validator language server
type Server struct {
	env *cel.Env
}

func NewServer() (*Server, error) {
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
	scanner := bufio.NewScanner(os.Stdin)

	// Read each line from stdin
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue // Skiping empty lines
		}

		diagnostics := s.validateCEL(input)

		if len(diagnostics) > 0 {
			for _, diag := range diagnostics {
				fmt.Fprintf(os.Stderr, "Error: %s\n", diag.Message)
			}
		} else {
			fmt.Println("Valid CEL expression")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return err
	}

	return nil
}

func (s *Server) validateCEL(content string) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	content = strings.TrimSpace(content)
	if content == "" {
		return diagnostics
	}

	ast, issues := s.env.Parse(content)

	if issues != nil && issues.Err() != nil {
		//Adding diagnostics for checking syntax errors
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
			//Adding diagnostics for checking syntax errors
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
