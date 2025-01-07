package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Define flags for timeout and command
	timeout := flag.Int("timeout", 10, "Timeout in seconds for the command to execute")
	flag.Parse()

	// Retrieve the command and its arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No command provided.")
		os.Exit(1)
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	// Create the command
	cmd := exec.Command(cmdName, cmdArgs...)

	// Pipe stdout and stderr directly to the main command's stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	exitChan := make(chan error, 1)

	go func() {
		exitChan <- cmd.Run()
	}()

	// Use a timer for the timeout
	select {
	case err := <-exitChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Println("Command executed successfully")
		}
	case <-time.After(time.Duration(*timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to kill process: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Command timed out")
		os.Exit(1)
	}
}
