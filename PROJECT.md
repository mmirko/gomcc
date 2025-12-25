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
