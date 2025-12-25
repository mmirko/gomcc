# Project Structure

## Core Files

### Source Code
- **main.go** - CLI argument parsing and main execution logic
- **config.go** - Configuration file loading, validation, and app management
- **executor.go** - App execution engine with dependency resolution

### Configuration Files
- **go.mod** - Go module definition
- **Makefile** - Build and installation automation

### Documentation
- **README.md** - Comprehensive user guide and documentation
- **QUICKSTART.md** - Quick start guide for getting up and running
- **LICENSE** - Project license

### Example Configurations
- **example-config.json** - Basic example from requirements (host A/B scenario)
- **example-webapp.json** - Complex web application stack example
- **example-test.json** - Test configuration with check apps

### Test Scripts
- **test-check-success.sh** - Sample successful check script
- **test-check-failure.sh** - Sample failing check script

### Build Artifacts
- **gomcc** - Compiled binary (created after build)
- **.gitignore** - Git ignore rules

## Key Features Implemented

### 1. Configuration Management
- JSON-based configuration file
- App type validation (check vs executable)
- Unique name validation
- Dependency validation
- Support for command arguments

### 2. App Types
- **Check Apps**: Exit-based status validation
- **Executable Apps**: Long-running processes

### 3. Dependency System
- Check-based dependencies
- `on_success` action execution
- `on_failure` action execution
- Dependency result caching

### 4. Tag System
- Tag-based filtering
- Multiple tag support
- Group-based launching

### 5. CLI Options
- `-f` - Configuration file path (default: ~/.gomcc.json)
- `-l` - List executable app names (one per line)
- `-L` - List all executable apps with detailed information
- `-t` - Filter by tags (comma-separated)
- `-v` - Verbose mode
- `-d` - Debug mode
- `-r` - Dry-run mode
- `-c` - Launch specific app
- `-g` - Launch all apps with tag
- `-e` - Test check app

### 6. Logging Levels
- Normal - Basic output
- Verbose - Detailed execution information
- Debug - Full diagnostic information

### 7. Process Management
- Process detachment for independence
- Continues running after gomcc exits
- Proper exit code handling

### 8. Additional Features
- List executable apps with resolved commands
- Default configuration file support (~/.gomcc.json)
- Check result caching
- Configuration validation
- Comprehensive error handling
- Execution summary
- Help system

## Build Commands

```bash
# Build
make build

# Clean
make clean

# Install (requires sudo)
make install

# Run example
make run-example

# Manual build
go build -o gomcc .
```

## Testing

All features have been tested:
✓ Configuration loading and validation
✓ Check app execution
✓ Executable app launching
✓ Dependency resolution
✓ Tag filtering
✓ Group launching
✓ Dry-run mode
✓ Verbose and debug modes
✓ Check result printing
✓ Command-line help

## Exit Codes

- 0 - All apps launched successfully
- 1 - One or more failures or configuration errors

## Project creation and prompts

This project was created by an AI language model based on detailed prompts describing the desired functionality, features, and structure. The AI generated the code, documentation, and examples iteratively to meet the specified requirements.

The Prompts are:

### Initial Prompt

Create a CLI linux launcher called gomcc written in Go with these features:
Load a JSON configuration file from a specified path.
The configuration file contains settings for every app to be launched.
The app can be of two types: the first type is a standard executable file snf its arguments, that spawns a process when launched;
The second type is a check app that exits when launched, returning a status code indicating success or failure.
The second type of app can be used as conditionals to determine whether to launch other executables.
The type of each app is specified in the configuration file for every app.
The apps can have dependencies on other apps, meaning that an app will only be launched if all its dependencies (checks) have succeeded.
Both types of apps are identified by a unique name in the configuration file.
Apps can be launched differently based on checks results of their dependencies. (there is not a single way to launch an app based on dependencies results)
The apps can also have tags associated with them, which can be used to filter which apps to launch and to group them.
The CLI tool should accept command line arguments to specify:

- The path to the JSON configuration file. (-f)
- A list of tags to filter which app to launch. (-t)
- An option to run in verbose mode, printing detailed logs of the execution process. (-v)
- A debug mode that prints even more detailed information for troubleshooting. (-d)
- An option to run in dry-run mode, where the tool only prints what would be executed without actually launching any apps. (-r)
- An to specify a single app to launch by its unique name. (-c appname)
- An option to launch all apps tagged with a specific tag. (-g tagname)
- A check printer to determine whether a check app exits with success or failure. (-e checkapp)
- Upon launched completion, the tool should exit with a status code indicating overall success (0) or failure (non-zero) based on the success of the launched apps

and their dependencies; The launched processes should not exit and continue running after the tool exits.

As example: an app called "startprogramtest" of type executable. We have two hosts A and B. The app "startprogramtest" is installed on host A, but reachable
from host B via ssh. Two check apps "checkonhostA" and "checkonhostB", gives success on host A and B respectively.
the app should be launched as "startprogramtest" only if "checkonhostA" cheecks succeeds, and and as "ssh hostB startprogramtest" if "checkonhostB" check succeeds.
so the configuration file should look like this:

```JSON
{
  "apps": [
    {
      "name": "checkonhostA",
      "type": "check",
      "command": "checkonhostAcommand",
      "tags": ["hostA"]
    },
    {
      "name": "checkonhostB",
      "type": "check",
      "command": "checkonhostBcommand",
      "tags": ["hostB"]
    },
    {
      "name": "startprogramtest",
      "type": "executable",
      "command": "startprogramtest",
      "dependencies": {
	"checkonhostA": {
	  "on_success": "startprogramtest"
	},
	"checkonhostB": {
	  "on_success": "ssh hostB startprogramtest"
	}
      },
      "tags": ["program", "test"]
    }
  ]
}
```

Along with the code, provide a README file that explains how to build and use the tool, with examples.

### Second Prompt

Can you upgrade the project:

- Create a new argument (-l) that lists the executables according to the current configuration of tags etc.
- make the .gomcc.json file in the user home directory the default if the -f flag is not specified
- when executing without any argument the command should not execute anything, instead it launches something. This is an error! look into it and fix it.

### Third Prompt

that ok! Just a little thing: use the -L for the list argument (the one currently implemented) and with the -l make e list of just the names of the executables one for each line, with no other information.
