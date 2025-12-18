# Quick Start Guide

## Building gomcc

```bash
# Using make
make build

# Or using go directly
go build -o gomcc .
```

## Installing gomcc

```bash
# Install to /usr/local/bin (requires sudo)
make install

# Or manually
sudo cp gomcc /usr/local/bin/
```

## Running the Examples

### 1. Test Example Configuration

```bash
# Dry-run with verbose output
./gomcc -f example-config.json -r -v

# Launch specific app in debug mode
./gomcc -f example-config.json -c startprogramtest -d -r
```

### 2. Web Application Example

```bash
# Launch all production apps (dry-run)
./gomcc -f example-webapp.json -t production -r -v

# Launch only backend services (dry-run)
./gomcc -f example-webapp.json -g backend -r -v

# Launch only development apps (dry-run)
./gomcc -f example-webapp.json -t development -r -v
```

### 3. Test Check Apps

```bash
# Test a successful check
./gomcc -f example-test.json -e check_success

# Test a failing check
./gomcc -f example-test.json -e check_failure

# Test conditional execution based on checks (dry-run)
./gomcc -f example-test.json -c conditional_app -r -d
```

## Creating Your Own Configuration

1. Create a JSON file with your app definitions
2. Define check apps for validations
3. Define executable apps with optional dependencies
4. Use tags for organization
5. Test with dry-run mode first: `./gomcc -f yourconfig.json -r -v`

## Common Usage Patterns

### Launch All Apps
```bash
gomcc -f config.json
```

### Filter by Tags
```bash
gomcc -f config.json -t tag1,tag2
```

### Launch Single App
```bash
gomcc -f config.json -c appname
```

### Launch All Apps with Tag
```bash
gomcc -f config.json -g tagname
```

### Test a Check
```bash
gomcc -f config.json -e checkname
```

### Debug Mode
```bash
gomcc -f config.json -d -r
```

## Tips

- Always test with `-r` (dry-run) first
- Use `-v` or `-d` for detailed logging
- Check apps are cached during execution
- Launched processes continue after gomcc exits
- Exit code 0 = success, 1 = failure
