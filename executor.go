package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// LogLevel represents the verbosity level
type LogLevel int

const (
	LogNormal LogLevel = iota
	LogVerbose
	LogDebug
)

// Executor handles the execution of apps
type Executor struct {
	config     *Config
	logLevel   LogLevel
	dryRun     bool
	checkCache map[string]bool // Cache for check results
}

// NewExecutor creates a new executor
func NewExecutor(config *Config, logLevel LogLevel, dryRun bool) *Executor {
	return &Executor{
		config:     config,
		logLevel:   logLevel,
		dryRun:     dryRun,
		checkCache: make(map[string]bool),
	}
}

// log prints a message at the specified level
func (e *Executor) log(level LogLevel, format string, args ...interface{}) {
	if e.logLevel >= level {
		fmt.Printf(format+"\n", args...)
	}
}

// ExecuteCheck runs a check app and returns true if successful
func (e *Executor) ExecuteCheck(app *App) (bool, error) {
	if app.Type != TypeCheck {
		return false, fmt.Errorf("app '%s' is not a check type", app.Name)
	}

	// Check cache first
	if result, exists := e.checkCache[app.Name]; exists {
		e.log(LogDebug, "[DEBUG] Using cached result for check '%s': %v", app.Name, result)
		return result, nil
	}

	e.log(LogVerbose, "[VERBOSE] Executing check: %s", app.Name)
	e.log(LogDebug, "[DEBUG] Check command: %s %v", app.Command, app.Args)

	if e.dryRun {
		e.log(LogNormal, "[DRY-RUN] Would execute check: %s %s", app.Command, strings.Join(app.Args, " "))
		// In dry-run mode, assume success
		e.checkCache[app.Name] = true
		return true, nil
	}

	cmd := exec.Command(app.Command, app.Args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	success := err == nil

	// Cache the result
	e.checkCache[app.Name] = success

	if success {
		e.log(LogVerbose, "[VERBOSE] Check '%s' succeeded", app.Name)
	} else {
		e.log(LogVerbose, "[VERBOSE] Check '%s' failed: %v", app.Name, err)
	}

	return success, nil
}

// PrintCheckResult executes a check and prints the result
func (e *Executor) PrintCheckResult(app *App) error {
	if app.Type != TypeCheck {
		return fmt.Errorf("app '%s' is not a check type", app.Name)
	}

	fmt.Printf("Executing check: %s\n", app.Name)
	fmt.Printf("Command: %s %s\n", app.Command, strings.Join(app.Args, " "))

	cmd := exec.Command(app.Command, app.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err == nil {
		fmt.Printf("Result: SUCCESS (exit code 0)\n")
		return nil
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			fmt.Printf("Result: FAILURE (exit code %d)\n", exitCode)
		} else {
			fmt.Printf("Result: FAILURE (%v)\n", err)
		}
		return err
	}
}

// ResolveCommand is the public version of resolveCommand
func (e *Executor) ResolveCommand(app *App) (string, []string, error) {
	return e.resolveCommand(app)
}

// resolveCommand determines the actual command to execute based on dependencies
func (e *Executor) resolveCommand(app *App) (string, []string, error) {
	if len(app.Dependencies) == 0 {
		// No dependencies, use the default command
		return app.Command, app.Args, nil
	}

	e.log(LogDebug, "[DEBUG] Resolving command for app '%s' with dependencies", app.Name)

	// Check all dependencies and build the command
	for depName, action := range app.Dependencies {
		depApp := e.config.GetApp(depName)
		if depApp == nil {
			return "", nil, fmt.Errorf("dependency '%s' not found", depName)
		}

		if depApp.Type != TypeCheck {
			return "", nil, fmt.Errorf("dependency '%s' is not a check type", depName)
		}

		success, err := e.ExecuteCheck(depApp)
		if err != nil {
			e.log(LogDebug, "[DEBUG] Error executing dependency check '%s': %v", depName, err)
		}

		if success && action.OnSuccess != "" {
			e.log(LogDebug, "[DEBUG] Dependency '%s' succeeded, using on_success command", depName)
			return e.parseCommand(action.OnSuccess)
		} else if !success && action.OnFailure != "" {
			e.log(LogDebug, "[DEBUG] Dependency '%s' failed, using on_failure command", depName)
			return e.parseCommand(action.OnFailure)
		}
	}

	// If no dependency action matched, use default command
	e.log(LogDebug, "[DEBUG] No dependency action matched, using default command")
	return app.Command, app.Args, nil
}

// parseCommand splits a command string into command and arguments
func (e *Executor) parseCommand(cmdStr string) (string, []string, error) {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("empty command string")
	}
	return parts[0], parts[1:], nil
}

// ExecuteApp launches an executable app
func (e *Executor) ExecuteApp(app *App) error {
	if app.Type != TypeExecutable {
		return fmt.Errorf("app '%s' is not an executable type", app.Name)
	}

	e.log(LogVerbose, "[VERBOSE] Preparing to execute app: %s", app.Name)

	// Resolve the actual command based on dependencies
	cmd, args, err := e.resolveCommand(app)
	if err != nil {
		return fmt.Errorf("failed to resolve command for app '%s': %w", app.Name, err)
	}

	fullCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	e.log(LogDebug, "[DEBUG] Resolved command: %s", fullCmd)

	if e.dryRun {
		e.log(LogNormal, "[DRY-RUN] Would execute app '%s': %s", app.Name, fullCmd)
		return nil
	}

	e.log(LogNormal, "Launching app '%s': %s", app.Name, fullCmd)

	// Create command
	execCmd := exec.Command(cmd, args...)

	// Detach the process so it continues running after we exit
	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the process
	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("failed to start app '%s': %w", app.Name, err)
	}

	e.log(LogVerbose, "[VERBOSE] Successfully launched app '%s' with PID %d", app.Name, execCmd.Process.Pid)

	// Don't wait for the process to finish - let it run independently
	return nil
}

// CanExecuteApp checks if an app can be executed based on its dependencies
func (e *Executor) CanExecuteApp(app *App) (bool, error) {
	if len(app.Dependencies) == 0 {
		return true, nil
	}

	e.log(LogDebug, "[DEBUG] Checking if app '%s' can be executed", app.Name)

	// For an app to be executable, at least one dependency action must be satisfied
	canExecute := false

	for depName, action := range app.Dependencies {
		depApp := e.config.GetApp(depName)
		if depApp == nil {
			return false, fmt.Errorf("dependency '%s' not found", depName)
		}

		if depApp.Type != TypeCheck {
			return false, fmt.Errorf("dependency '%s' is not a check type", depName)
		}

		success, err := e.ExecuteCheck(depApp)
		if err != nil {
			e.log(LogDebug, "[DEBUG] Error executing dependency check '%s': %v", depName, err)
		}

		// If there's an action for this result, the app can be executed
		if (success && action.OnSuccess != "") || (!success && action.OnFailure != "") {
			canExecute = true
			e.log(LogDebug, "[DEBUG] Dependency '%s' allows execution (success=%v)", depName, success)
		}
	}

	if !canExecute {
		e.log(LogDebug, "[DEBUG] No dependency satisfied for app '%s'", app.Name)
	}

	return canExecute, nil
}
