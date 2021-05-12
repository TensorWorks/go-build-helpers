package process

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executes the specified command and waits for it to complete
func Run(command []string, dir *string, env *map[string]string) (error) {
	
	// Print the command prior to running it
	log.Print(command)
	
	// If the specified executable is not a filesystem path then perform PATH lookup
	executable := command[0]
	if filepath.Base(executable) == executable {
		if resolved, err := exec.LookPath(executable); err != nil {
			return err
		} else {
			executable = resolved
		}
	}
	
	// If no working directory was specified then use the current working directory
	workingDir := ""
	if dir != nil {
		workingDir = *dir
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		} else {
			workingDir = cwd
		}
	}
	
	// Merge any supplied environment variables with those of the parent process
	environment := os.Environ()
	if env != nil {
		
		// Gather the list of keys for our supplied environment variables
		keys := []string{}
		for key := range *env {
			keys = append(keys, key)
		}
		
		// Prevent duplicate entries by stripping out any existing values for our keys
		// (This ensures the underlying operating system doesn't ignore our supplied values as a result of them appearing after the system defaults)
		environment = stripKeys(environment, keys)
		
		// Append the supplied environment variables
		for k, v := range *env {
			environment = append(environment, fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	// Attempt to start the child process, ensuring it inherits its stdout and stderr streams from the parent
	process, err := os.StartProcess(executable, command, &os.ProcAttr{
		Dir: workingDir,
		Env: environment,
		Files: []*os.File{nil, os.Stdout, os.Stderr},
	})
	
	// Verify that the child process started successfully
	if err != nil {
		return err
	}
	
	// Wait for the child process to complete
	status, err := process.Wait()
	if err != nil {
		return err
	}
	
	// Verify that the child process completed successfully
	if status.Success() == false {
		return errors.New(fmt.Sprint("Command ", command, " terminated with exit code ", status.ExitCode()))
	}
	
	return nil
}

// Strips the specified list of keys from an array of environment variables
func stripKeys(env []string, keys []string) []string {
	
	// Construct a map of our keys for faster searching
	keysMap := map[string]bool{}
	for _, key := range keys {
		keysMap[key] = true
	}
	
	// Strip out any environment variable entries matching our keys
	stripped := []string{}
	for _, entry := range env {
		components := strings.SplitN(entry, "=", 2)
		if len(components) == 2 && !keysMap[components[0]] {
			stripped = append(stripped, entry)
		}
	}
	
	return stripped
}
