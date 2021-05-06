package module

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
