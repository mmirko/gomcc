# gomcc - Go Multi-Component Commander

A flexible CLI launcher for managing application dependencies and conditional execution based on check results.

## Features

- **Two Types of Apps:**
  - **Executable Apps**: Standard programs that spawn processes and continue running
  - **Check Apps**: Validation scripts that exit with status codes (0 for success, non-zero for failure)

- **Dependency-Based Execution**: Apps can have dependencies on check apps, with different commands executed based on check results

- **Tag-Based Filtering**: Organize and filter apps using tags

- **Multiple Execution Modes**:
  - List all executable apps
  - Launch all apps
  - Launch apps by tag(s)
  - Launch a specific app
  - Test a check app

- **Logging Levels**: Normal, verbose, and debug modes for different levels of detail

- **Dry-Run Mode**: Preview what would be executed without actually running anything

- **Process Independence**: Launched processes continue running after gomcc exits

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mmirko/gomcc.git
cd gomcc

# Build the binary
go build -o gomcc

# Optional: Install to your PATH
sudo mv gomcc /usr/local/bin/
```

## Configuration File

The configuration file is a JSON file that defines all apps and their dependencies.

### Configuration Structure

```json
{
  "apps": [
    {
      "name": "app_name",
      "type": "executable|check",
      "command": "command_to_execute",
      "args": ["arg1", "arg2"],
      "tags": ["tag1", "tag2"],
      "dependencies": {
        "check_app_name": {
          "on_success": "command to run if check succeeds",
          "on_failure": "command to run if check fails"
        }
      }
    }
  ]
}
```

### Field Descriptions

- **name** (required): Unique identifier for the app
- **type** (required): Either `"check"` or `"executable"`
- **command** (required): The command to execute
- **args** (optional): Array of command arguments
- **tags** (optional): Array of tags for filtering and grouping
- **dependencies** (optional): Map of check app names to dependency actions
  - **on_success**: Command to execute if the check succeeds
  - **on_failure**: Command to execute if the check fails

## Usage

### Command-Line Options

```
-f <path>       Path to the JSON configuration file (default: ~/.gomcc.json)
-l              List executable app names (one per line)
-L              List all executable apps with detailed information
-t <tags>       Comma-separated list of tags to filter apps
-v              Enable verbose mode
-d              Enable debug mode (implies verbose)
-r              Enable dry-run mode (don't actually execute)
-c <appname>    Launch a specific app by name
-g <tagname>    Launch all apps with a specific tag
-e <checkapp>   Execute and print result of a check app
```

### Examples

#### 1. List Executable Apps

```bash
# List app names only (one per line)
gomcc -l

# List with detailed information
gomcc -L

# List from specific config
gomcc -f config.json -l

# List with details from specific config
gomcc -f config.json -L

# List apps filtered by tags
gomcc -l -t web,backend
```

#### 2. Launch All Apps

```bash
# Using default config (~/.gomcc.json)
gomcc

# Using specific config
gomcc -f config.json
```

#### 3. Launch Apps with Specific Tags

```bash
# Launch apps tagged with 'web' or 'backend'
gomcc -f config.json -t web,backend
```

#### 4. Launch a Specific App

```bash
gomcc -f config.json -c myapp
```

#### 5. Launch All Apps with a Specific Tag

```bash
gomcc -f config.json -g production
```

#### 6. Test a Check App

```bash
gomcc -f config.json -e checkapp
```

#### 7. Dry-Run Mode

```bash
# See what would be executed without running anything
gomcc -f config.json -r
```

#### 8. Verbose Mode

```bash
gomcc -f config.json -v
```

#### 9. Debug Mode

```bash
gomcc -f config.json -d
```

## Complete Example

This example demonstrates a program that can run on different hosts based on check results.

### Configuration File (config.json)

```json
{
  "apps": [
    {
      "name": "checkonhostA",
      "type": "check",
      "command": "ping",
      "args": ["-c", "1", "hostA"],
      "tags": ["hostA"]
    },
    {
      "name": "checkonhostB",
      "type": "check",
      "command": "ping",
      "args": ["-c", "1", "hostB"],
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

### How It Works

1. **checkonhostA**: Pings hostA to check if it's reachable
2. **checkonhostB**: Pings hostB to check if it's reachable
3. **startprogramtest**: 
   - If `checkonhostA` succeeds, runs `startprogramtest` locally
   - If `checkonhostB` succeeds, runs `ssh hostB startprogramtest` remotely
   - Can be launched both ways if both checks succeed

### Running the Example

```bash
# Launch the program
gomcc -f config.json -c startprogramtest

# Dry-run to see what would happen
gomcc -f config.json -c startprogramtest -r -v

# Test individual checks
gomcc -f config.json -e checkonhostA
gomcc -f config.json -e checkonhostB
```

## Advanced Example: Web Application Stack

```json
{
  "apps": [
    {
      "name": "check_docker",
      "type": "check",
      "command": "docker",
      "args": ["info"],
      "tags": ["infrastructure"]
    },
    {
      "name": "check_port_80",
      "type": "check",
      "command": "sh",
      "args": ["-c", "! netstat -tuln | grep ':80 '"],
      "tags": ["infrastructure"]
    },
    {
      "name": "database",
      "type": "executable",
      "command": "docker",
      "args": ["run", "-d", "--name", "mydb", "postgres:latest"],
      "dependencies": {
        "check_docker": {
          "on_success": "docker run -d --name mydb postgres:latest"
        }
      },
      "tags": ["backend", "database"]
    },
    {
      "name": "webserver",
      "type": "executable",
      "command": "nginx",
      "dependencies": {
        "check_port_80": {
          "on_success": "nginx"
        }
      },
      "tags": ["web", "frontend"]
    },
    {
      "name": "api_server",
      "type": "executable",
      "command": "/opt/myapp/api-server",
      "args": ["--port", "3000"],
      "tags": ["backend", "api"]
    }
  ]
}
```

### Running Different Scenarios

```bash
# Launch only backend components
gomcc -f webapp.json -t backend

# Launch only frontend
gomcc -f webapp.json -t frontend

# Launch everything
gomcc -f webapp.json

# Check if Docker is available
gomcc -f webapp.json -e check_docker

# Dry-run the entire stack
gomcc -f webapp.json -r -v
```

## Exit Codes

- **0**: All apps launched successfully
- **1**: One or more apps failed to launch or configuration error

## Tips and Best Practices

1. **Use Specific Checks**: Create check apps that validate exact prerequisites for your executables

2. **Tag Wisely**: Use tags to create logical groups (e.g., "production", "development", "backend", "frontend")

3. **Test Checks First**: Use the `-e` flag to test individual check apps before running the full configuration

4. **Dry-Run Before Production**: Always use `-r -v` to preview execution in complex scenarios

5. **Cache Check Results**: gomcc automatically caches check results within a single execution to avoid redundant checks

6. **Process Independence**: Launched processes are detached and will continue running after gomcc exits

7. **Dependency Actions**: You can specify both `on_success` and `on_failure` actions for the same dependency to handle both scenarios

## Troubleshooting

### App Not Launching

```bash
# Use debug mode to see detailed execution flow
gomcc -f config.json -c myapp -d
```

### Check Dependencies

```bash
# Test individual checks
gomcc -f config.json -e mycheck
```

### Validate Configuration

```bash
# Dry-run with debug mode
gomcc -f config.json -r -d
```

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
