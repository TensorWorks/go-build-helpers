package module

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tensorworks/go-build-helpers/pkg/process"
)

// Returns the absolute filesystem path to the directory which will hold any tools needed for codegen
func (module *Module) CodegenToolsDir() (string) {
	return filepath.Join(module.BuildDir(), "tools", runtime.GOOS, runtime.GOARCH)
}

// Installs the specified Go tools to our codegen tools directory
func (module *Module) InstallGoTools(tools []string) (error) {
	
	// Set the GOBIN environment variable to our codegen tools directory
	env := &map[string]string{
		"GOBIN": module.CodegenToolsDir(),
	}
	
	// Iterate over the list of Go tools and install each in turn
	for _, tool := range tools {
		if err := process.Run([]string{"go", "install", tool}, &module.RootDir, env); err != nil {
			return err
		}
	}
	
	// If we reached this point then all tools were installed successfully
	return nil
}

// Performs code generation for the Go module
func (module *Module) Generate() (error) {
	
	// Invoke `go generate` with the our codegen tools directory appended to the PATH
	return process.Run(
		[]string{"go", "generate", "./..."},
		&module.RootDir,
		&map[string]string{
			"PATH": fmt.Sprint(os.Getenv("PATH"), string(os.PathListSeparator), module.CodegenToolsDir()),
		},
	)
	
}
