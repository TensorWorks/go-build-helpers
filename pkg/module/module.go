package module

import (
	"os"
)

// Contains information about a Go module
type Module struct {
	
	// The absolute filesystem path to the root of the module
	RootDir string
	
}

// Creates a Module object for the Go module in the current working directory
func ModuleInCwd() (*Module, error) {
	
	// Attempt to resolve the path to the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	
	// Create a Codegen object for the current working directory
	return &Module{RootDir: cwd}, nil
}
