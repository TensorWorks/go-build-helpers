package module

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tensorworks/go-build-helpers/pkg/process"
)

// Callers can use this constant when calling BuildBinariesFor[...]() to specify that the module's binaries directory should be used
const DefaultBinDir = ""

// Represents a matrix of build configurations
type BuildMatrix struct {
	Platforms []string
	Architectures []string
	Ignore []string
}

// Returns the absolute filesystem path to the directory which will hold the binaries for the module
func (module *Module) BinariesDir() (string) {
	return filepath.Join(module.RootDir, "bin")
}

// Builds all of the executables in the Go module for the specified build context
func (module *Module) BuildBinariesForContext(binDir string, options BuildOptions, platform string, architecture string) (error) {
	
	// If either of our options arrays are nil then set them to empty arrays
	if (options.BuildTags == nil) {
		options.BuildTags = []string{}
	}
	if (options.AdditionalFlags == nil) {
		options.AdditionalFlags = []string{}
	}
	
	// If no output directory was specified then use the binaries directory for the module
	if binDir == DefaultBinDir {
		binDir = module.BinariesDir()
	}
	
	// If we are using a directory prefix naming scheme then append the context directories to the output path
	if options.Scheme == PrefixedDirs {
		binDir = filepath.Join(binDir, platform, architecture)
	}
	
	// Ensure the output directory exists
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return err
	}
	
	// If we are using a filename suffix then we place the binaries in a temporary staging directory prior to renaming them
	// (This allows us to avoid clobbering any existing binaries produced using the undecorated naming scheme)
	origBinDir := binDir
	if options.Scheme == SuffixedFilenames {
		binDir = filepath.Join(module.BuildDir(), "staging", platform, architecture)
	}
	
	// Prepare the flags for invoking `go build`
	flags := []string{"go", "build", "-o", fmt.Sprint(binDir, string(os.PathSeparator)), "-tags", strings.Join(options.BuildTags, ",")}
	flags = append(flags, options.AdditionalFlags...)
	flags = append(flags, "./...")
	
	// Invoke `go build` with the appropriate flags and environment variables
	// (Note: we append a trailing slash to the output directory to ensure `go build` always interprets it as a directory)
	err := process.Run(
		flags,
		&module.RootDir,
		&map[string]string{
			"GOOS": platform,
			"GOARCH": architecture,
		},
	)
	
	// If we are using a filename suffix then move the binaries from the temporary directory to the proper output directory
	if options.Scheme == SuffixedFilenames {
		
		// Retrieve the list of binary files
		binaries, err := filepath.Glob(filepath.Join(binDir, "*"))
		if err != nil {
			return err
		}
		
		// Move each binary file to the proper output directory, adding the filename suffix for the build context
		for _, binary := range binaries {
			
			// Isolate the components of the binary's filename
			filename := filepath.Base(binary)
			extension := filepath.Ext(filename)
			executable := strings.TrimSuffix(filename, extension)
			
			// Construct the final filesystem path for the binary and move it into place
			newpath := filepath.Join(origBinDir, fmt.Sprint(executable, "-", platform, "-", architecture, extension))
			if err := os.Rename(binary, newpath); err != nil {
				return err
			}
		}
	}
	
	return err
}

// Builds all of the executables in the Go module for the host system's build context
func (module *Module) BuildBinariesForHost(binDir string, options BuildOptions) (error) {
	return module.BuildBinariesForContext(binDir, options, runtime.GOOS, runtime.GOARCH)
}

// Builds all of the executables in the Go module for all specified build context combinations
func (module *Module) BuildBinariesForMatrix(binDir string, options BuildOptions, matrix BuildMatrix) (error) {
	
	// Verify that the user did not specify the undecorated naming scheme, which would result in binaries being clobbered
	if options.Scheme == Undecorated {
		return errors.New("using the undecorated naming scheme when building a matrix of configurations would clobber binaries")
	}
	
	// Construct a map of our ignored build context combinations for faster searching
	ignoreMap := map[string]bool{}
	for _, key := range matrix.Ignore {
		ignoreMap[key] = true
	}
	
	// Iterate through the build context combinations and build each one
	for _, platform := range matrix.Platforms {
		for _, architecture := range matrix.Architectures {
			
			// If the current build context combination is ignored then skip building it
			if _, ignored := ignoreMap[fmt.Sprint(platform, "/", architecture)]; ignored {
				continue
			}
			
			if err := module.BuildBinariesForContext(binDir, options, platform, architecture); err != nil {
				return err
			}
		}
	}
	
	// If we reached this point then all combinations built successfully
	return nil
}
