package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CLI represents the command-line interface
type CLI struct {
	configPath        string
	tags              string
	verbose           bool
	debug             bool
	dryRun            bool
	appName           string
	groupTag          string
	checkApp          string
	listApps          bool
	listAppsDetailed  bool
}

// ParseArgs parses command-line arguments
func (c *CLI) ParseArgs() {
	// Get default config path
	homeDir, err := os.UserHomeDir()
	defaultConfig := ""
	if err == nil {
		defaultConfig = filepath.Join(homeDir, ".gomcc.json")
	}
	
	flag.StringVar(&c.configPath, "f", defaultConfig, "Path to the JSON configuration file")
	flag.StringVar(&c.tags, "t", "", "Comma-separated list of tags to filter apps")
	flag.BoolVar(&c.verbose, "v", false, "Enable verbose mode")
	flag.BoolVar(&c.debug, "d", false, "Enable debug mode (implies verbose)")
	flag.BoolVar(&c.dryRun, "r", false, "Enable dry-run mode (don't actually execute)")
	flag.StringVar(&c.appName, "c", "", "Launch a specific app by name")
	flag.StringVar(&c.groupTag, "g", "", "Launch all apps with a specific tag")
	flag.StringVar(&c.checkApp, "e", "", "Execute and print result of a check app")
	flag.BoolVar(&c.listApps, "l", false, "List executable app names (one per line)")
	flag.BoolVar(&c.listAppsDetailed, "L", false, "List all executable apps with detailed information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "gomcc - A flexible CLI launcher for managing application dependencies\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                   # Launch all apps (uses ~/.gomcc.json)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -l                                # List executable app names\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -L                                # List executable apps with details\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f config.json                    # Launch all apps with specific config\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t web,backend                    # Launch apps with web or backend tags\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c myapp                          # Launch only 'myapp'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -g production                     # Launch all apps tagged 'production'\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -e checkapp                       # Test a check app\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -v                                # Launch with verbose logging\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -r                                # Dry-run mode\n", os.Args[0])
	}

	flag.Parse()

	// Check if any action is specified
	hasAction := c.listApps || c.listAppsDetailed || c.checkApp != "" || c.appName != "" || c.groupTag != "" || len(flag.Args()) > 0
	
	// If no action specified and no other arguments, show usage
	if !hasAction && flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Validate config file path
	if c.configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: config file path not specified and no default found\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Debug mode implies verbose
	if c.debug {
		c.verbose = true
	}
}

// GetLogLevel returns the appropriate log level based on flags
func (c *CLI) GetLogLevel() LogLevel {
	if c.debug {
		return LogDebug
	}
	if c.verbose {
		return LogVerbose
	}
	return LogNormal
}

// GetTagsList returns the list of tags as a slice
func (c *CLI) GetTagsList() []string {
	if c.tags == "" {
		return nil
	}
	parts := strings.Split(c.tags, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// equalStringSlices compares two string slices for equality
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	cli := &CLI{}
	cli.ParseArgs()

	// Load configuration
	config, err := LoadConfig(cli.configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create executor
	executor := NewExecutor(config, cli.GetLogLevel(), cli.dryRun)

	// Handle list apps mode (simple - just names)
	if cli.listApps {
		tags := cli.GetTagsList()
		var appsToList []App
		
		if cli.appName != "" {
			// List specific app
			app := config.GetApp(cli.appName)
			if app == nil {
				fmt.Fprintf(os.Stderr, "Error: app '%s' not found\n", cli.appName)
				os.Exit(1)
			}
			appsToList = []App{*app}
		} else if cli.groupTag != "" {
			// List apps with specific tag
			appsToList = config.GetAppsByTag(cli.groupTag)
		} else if len(tags) > 0 {
			// List apps filtered by tags
			appsToList = config.GetAppsByTags(tags)
		} else {
			// List all apps
			appsToList = config.Apps
		}
		
		// Print only executable names, one per line
		for _, app := range appsToList {
			if app.Type == TypeExecutable {
				fmt.Println(app.Name)
			}
		}
		os.Exit(0)
	}

	// Handle list apps mode (detailed)
	if cli.listAppsDetailed {
		tags := cli.GetTagsList()
		var appsToList []App
		
		if cli.appName != "" {
			// List specific app
			app := config.GetApp(cli.appName)
			if app == nil {
				fmt.Fprintf(os.Stderr, "Error: app '%s' not found\n", cli.appName)
				os.Exit(1)
			}
			appsToList = []App{*app}
		} else if cli.groupTag != "" {
			// List apps with specific tag
			appsToList = config.GetAppsByTag(cli.groupTag)
		} else if len(tags) > 0 {
			// List apps filtered by tags
			appsToList = config.GetAppsByTags(tags)
		} else {
			// List all apps
			appsToList = config.Apps
		}
		
		fmt.Println("Executable Apps:")
		fmt.Println("================")
		count := 0
		for _, app := range appsToList {
			if app.Type == TypeExecutable {
				count++
				fmt.Printf("\nName: %s\n", app.Name)
				if len(app.Tags) > 0 {
					fmt.Printf("  Tags: %s\n", strings.Join(app.Tags, ", "))
				}
				
				// Resolve command based on dependencies
				cmd, args, err := executor.ResolveCommand(&app)
				if err != nil {
					fmt.Printf("  Command: %s %s (error: %v)\n", app.Command, strings.Join(app.Args, " "), err)
				} else {
					if cmd != app.Command || !equalStringSlices(args, app.Args) {
						fmt.Printf("  Default Command: %s %s\n", app.Command, strings.Join(app.Args, " "))
						fmt.Printf("  Resolved Command: %s %s\n", cmd, strings.Join(args, " "))
					} else {
						fmt.Printf("  Command: %s %s\n", cmd, strings.Join(args, " "))
					}
				}
				
				if len(app.Dependencies) > 0 {
					fmt.Printf("  Dependencies:\n")
					for depName, action := range app.Dependencies {
						fmt.Printf("    - %s\n", depName)
						if action.OnSuccess != "" {
							fmt.Printf("        on_success: %s\n", action.OnSuccess)
						}
						if action.OnFailure != "" {
							fmt.Printf("        on_failure: %s\n", action.OnFailure)
						}
					}
				}
			}
		}
		
		if count == 0 {
			fmt.Println("\nNo executable apps found.")
		} else {
			fmt.Printf("\nTotal: %d executable app(s)\n", count)
		}
		os.Exit(0)
	}

	// Handle check app mode
	if cli.checkApp != "" {
		app := config.GetApp(cli.checkApp)
		if app == nil {
			fmt.Fprintf(os.Stderr, "Error: app '%s' not found\n", cli.checkApp)
			os.Exit(1)
		}
		if app.Type != TypeCheck {
			fmt.Fprintf(os.Stderr, "Error: app '%s' is not a check type\n", cli.checkApp)
			os.Exit(1)
		}

		err := executor.PrintCheckResult(app)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	var appsToLaunch []App
	var executionMode string

	// Determine which apps to launch
	if cli.appName != "" {
		// Launch specific app
		executionMode = fmt.Sprintf("app '%s'", cli.appName)
		app := config.GetApp(cli.appName)
		if app == nil {
			fmt.Fprintf(os.Stderr, "Error: app '%s' not found\n", cli.appName)
			os.Exit(1)
		}
		appsToLaunch = []App{*app}
	} else if cli.groupTag != "" {
		// Launch all apps with a specific tag
		executionMode = fmt.Sprintf("apps with tag '%s'", cli.groupTag)
		appsToLaunch = config.GetAppsByTag(cli.groupTag)
		if len(appsToLaunch) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: no apps found with tag '%s'\n", cli.groupTag)
		}
	} else {
		// Launch all apps (filtered by tags if specified)
		tags := cli.GetTagsList()
		if len(tags) > 0 {
			executionMode = fmt.Sprintf("apps with tags [%s]", strings.Join(tags, ", "))
			appsToLaunch = config.GetAppsByTags(tags)
		} else {
			executionMode = "all apps"
			appsToLaunch = config.Apps
		}
	}

	executor.log(LogVerbose, "[VERBOSE] Execution mode: %s", executionMode)
	executor.log(LogVerbose, "[VERBOSE] Found %d app(s) to process", len(appsToLaunch))

	// Execute apps
	successCount := 0
	failureCount := 0
	skippedCount := 0

	for _, app := range appsToLaunch {
		// Skip check apps in normal execution
		if app.Type == TypeCheck {
			executor.log(LogDebug, "[DEBUG] Skipping check app '%s' in execution", app.Name)
			skippedCount++
			continue
		}

		// Check if app can be executed based on dependencies
		canExecute, err := executor.CanExecuteApp(&app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking dependencies for app '%s': %v\n", app.Name, err)
			failureCount++
			continue
		}

		if !canExecute {
			executor.log(LogVerbose, "[VERBOSE] Skipping app '%s' - dependencies not satisfied", app.Name)
			skippedCount++
			continue
		}

		// Execute the app
		err = executor.ExecuteApp(&app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing app '%s': %v\n", app.Name, err)
			failureCount++
		} else {
			successCount++
		}
	}

	// Print summary
	if cli.verbose || cli.dryRun {
		fmt.Println()
		fmt.Printf("Execution Summary:\n")
		fmt.Printf("  Successfully launched: %d\n", successCount)
		fmt.Printf("  Failed to launch:      %d\n", failureCount)
		fmt.Printf("  Skipped:               %d\n", skippedCount)
	}

	// Exit with appropriate status code
	if failureCount > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
