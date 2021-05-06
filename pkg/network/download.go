package network

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	filesystem "github.com/tensorworks/go-build-helpers/pkg/filesystem"
	module "github.com/tensorworks/go-build-helpers/pkg/module"
)

// Callers can use this constant for the Filename field of DownloadedFile to specify that the output filename should be auto-detected from the URL
const AutoDetectFilename = ""

// Represents a file that is downloaded from a URL
type DownloadedFile struct {
	
	// The URL from which the file is downloaded
	URL string
	
	// The filename used to store the file contents on the local filesystem
	Filename string
	
	// The Go module (if any) for which the file will be downloaded
	Module *module.Module
	
}

// Returns the absolute filesystem path for the local copy of the file
func (file *DownloadedFile) LocalPath() (string) {
	
	// If we are auto-detecting the filename from the URL then do so
	filename := file.Filename
	if filename == AutoDetectFilename {
		filename = path.Base(file.URL)
	}
	
	// If the file is for a module, place it in the module's downloads directory
	if file.Module != nil {
		return filepath.Join(file.Module.DownloadsDir(), filepath.Base(filename))
	}
	
	// Attempt to resolve the absolute path for the specified filename
	if abs, err := filepath.Abs(filename); err == nil {
		return abs
	} else {
		return filename
	}
}

// Determines if the file has already been downloaded
func (file *DownloadedFile) Exists() (bool) {
	return filesystem.Exists(file.LocalPath())
}

// Downloads the file
func (file *DownloadedFile) Download() (error) {
	
	// Resolve the absolute path to the output location on the local filesystem
	downloaded := file.LocalPath()
	
	// Ensure the parent directory of the output path exists
	if err := os.MkdirAll(filepath.Dir(downloaded), os.ModePerm); err != nil {
		return err
	}
	
	// Attempt to create the output file
	stream, err := os.Create(downloaded)
	if err != nil {
		return err
	}
	
	// Ensure the output file is closed when the function completes
	defer stream.Close()
	
	// Attempt to perform a HTTP GET request for the download URL
	response, err := http.Get(file.URL)
	if err != nil {
		return err
	}
	
	// Ensure the HTTP stream is closed when the function completes
	defer response.Body.Close()
	
	// Stream the data from the HTTP response to our output file
	if _, err := io.Copy(stream, response.Body); err != nil {
		return err
	}
	
	return nil
}
