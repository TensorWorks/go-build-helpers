package module

import (
	"path/filepath"
	"runtime"
)

// Returns the absolute filesystem path to the directory which will hold build-related files for the module
func (module *Module) BuildDir() (string) {
	return filepath.Join(module.RootDir, ".build")
}

// Returns the absolute filesystem path to the directory which will hold any downloaded files needed for the build process
func (module *Module) DownloadsDir() (string) {
	return filepath.Join(module.BuildDir(), "downloads", runtime.GOOS, runtime.GOARCH)
}
