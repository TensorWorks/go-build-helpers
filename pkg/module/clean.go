package module

import (
	"os"
)

// Removes all filesystem directories used by the build process for the module
func (module *Module) CleanAll() (error) {
	
	// Attempt to clean the binaries directory
	if err := module.CleanBinariesDir(); err != nil {
		return err
	}
	
	// Attempt to clean the build directory
	if err := module.CleanBuildDir(); err != nil {
		return err
	}
	
	return nil
}

// Removes the filesystem directory which holds compiled binaries for the module
func (module *Module) CleanBinariesDir() (error) {
	return os.RemoveAll(module.BinariesDir())
}

// Removes the filesystem directory which holds build-related files for the module
func (module *Module) CleanBuildDir() (error) {
	return os.RemoveAll(module.BuildDir())
}
