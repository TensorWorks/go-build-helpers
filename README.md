# Go Build Helpers

**Note: this module requires Go 1.16 or newer.**

The packages in this module are designed to reduce boilerplate code when writing build scripts in Go, whether they be [Magefiles](https://magefile.org/) or just plain Go files run with `go run`.

The following packages are provided:

- [filesystem](./pkg/filesystem): provides functionality related to filesystem operations not concisely covered by the Go standard library.

- [module](./pkg/module): provides functionality for performing code generation and compilation of Go modules. Most of the other helper packages are designed to work with the [Module](./pkg/module/module.go) type from this package.

- [network](./pkg/network): provides functionality for interacting with remote servers and downloading files.

- [process](./pkg/process): provides functionality for creating and interacting with child processes.

- [system](./pkg/system): provides functionality and constants related to the underlying operating system.

- [tools](./pkg/tools):
  
  - tools/[protoc](./pkg/tools/protoc): provides functionality for working with the [Google protocol buffers](https://developers.google.com/protocol-buffers) compiler.

- [validation](./pkg/validation): provides functionality for validating build step results and working with errors.


## Usage examples

### Failing when errors are encountered

```go
// +build never

package main

import (
	module "github.com/tensorworks/go-build-helpers/pkg/module"
	validation "github.com/tensorworks/go-build-helpers/pkg/validation"
)

func main() {
	
	// Create a build helper for the Go module in the current working directory
	mod, err := module.ModuleInCwd()
	
	// If an error occurred then log it and exit immediately
	validation.ExitIfError(err)
	
	// Do stuff with the module.Module object here
	// ...
}
```

### Performing code generation

```go
// +build never

package main

import (
	module "github.com/tensorworks/go-build-helpers/pkg/module"
	validation "github.com/tensorworks/go-build-helpers/pkg/validation"
)

func main() {
	
	// Create a build helper for the Go module in the current working directory
	mod, err := module.ModuleInCwd()
	validation.ExitIfError(err)
	
	// Install any Go tools that we require for code generation into our codegen tools directory
	err = mod.InstallGoTools([]string{
		"golang.org/x/tools/cmd/stringer@v0.1.0",
	})
	validation.ExitIfError(err)
	
	// Run `go generate` with our codegen tools directory appended to the PATH
	err = mod.Generate()
	validation.ExitIfError(err)
}
```

### Building executables

```go
// +build never

package main

import (
	module "github.com/tensorworks/go-build-helpers/pkg/module"
	validation "github.com/tensorworks/go-build-helpers/pkg/validation"
)

func main() {
	
	// Create a build helper for the Go module in the current working directory
	mod, err := module.ModuleInCwd()
	validation.ExitIfError(err)
	
	// Build binaries for any executables in the module (./cmd/XXX) for the host GOOS/GOARCH and place them in ./bin
	err = mod.BuildBinariesForHost(module.DefaultBinDir, module.Undecorated)
	validation.ExitIfError(err)
	
	// Alternatively, build binaries for executables using a matrix of GOOS/GOARCH combinations
	err = mod.BuildBinariesForMatrix(
		
		module.DefaultBinDir,
		
		// Note: this produces binaries with suffixed filenames ("cmd/mytool" becomes "bin/mytool-${GOOS}-${GOARCH}${GOEXE}")
		// If you want binaries in subdirectories instead ("cmd/mytool" becomes "bin/${GOOS}/${GOARCH}/mytool${GOEXE}") then use module.PrefixedDirs
		module.SuffixedFilenames,
		
		module.BuildMatrix{
			
			// This will build binaries for the following GOOS/GOARCH combinations:
			// - darwin/amd64
			// - linux/386
			// - linux/amd64
			// - windows/386
			// - windows/amd64
			
			Platforms: []string{"darwin", "linux", "windows"},
			Architectures: []string{"386", "amd64"},
			Ignore: []string{"darwin/386"},
		},
	)
	validation.ExitIfError(err)
}
```


## Legal

Copyright &copy; 2021, TensorWorks Pty Ltd. Licensed under the MIT License, see the file [LICENSE](./LICENSE) for details.
