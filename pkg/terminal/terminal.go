package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetEnvAsMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}

func StartSubshell() {
	// Define the subshell command (e.g., /bin/bash or /bin/sh)
	env := GetEnvAsMap()
	shell := env["SHELL"] //"/bin/bash"

	// Create the command to start the subshell
	fmt.Println("Session started with Project: project and Environment: env \ntype 'exit' to exit the session at any time")
	cmd := exec.Command(shell)

	// Set up the new environment variables
	newEnv := os.Environ()
	newEnv = append(newEnv, "MY_VAR=HelloWorld")
	cmd.Env = newEnv

	// Redirect standard input, output, and error to the subshell
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the subshell
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting subshell: %v\n", err)
		return
	}

	// Wait for the subshell to exit
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Subshell exited with error: %v\n", err)
	} else {
		fmt.Println("Subshell exited successfully")
	}
}
