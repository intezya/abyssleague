# Running Tests in the Websocket Service

This document provides instructions on how to run tests in the websocket service.

## Direct Commands

### Run All Tests

To run all tests in the project:

```bash
go test ./...
```

### Run Tests for a Specific Package

To run tests for a specific package:

```bash
go test websocket/internal/adapters/config
```

### Run Tests with Verbose Output

To run tests with verbose output:

```bash
go test -v ./...
```

### Run a Specific Test

To run a specific test:

```bash
go test -run TestIsDevMode websocket/internal/adapters/config
```

### Run Tests with Coverage

To run tests with coverage information:

```bash
go test -cover ./...
```

### Generate Coverage Profile

To generate a coverage profile:

```bash
go test -coverprofile=coverage.out ./...
```

To view the coverage report in HTML format:

```bash
go tool cover -html=coverage.out -o coverage.html
```

## Using the Scripts

Two scripts are provided to simplify running tests with various options:

- `run_tests.ps1` (PowerShell script)
- `run_tests.bat` (Batch file)

Both scripts support the same options and functionality.

### Basic Usage

Run all tests:

```powershell
# PowerShell
.\run_tests.ps1

# Batch
run_tests.bat
```

### Options

- `-v`: Verbose output
- `-cover`: Show coverage
- `-coverprofile <file>`: Generate coverage profile
- `-package <package>`: Run tests for a specific package
- `-test <test>`: Run a specific test

### Examples

Run tests with verbose output:

```powershell
# PowerShell
.\run_tests.ps1 -v

# Batch
run_tests.bat -v
```

Run tests for a specific package:

```powershell
# PowerShell
.\run_tests.ps1 -package "websocket/internal/adapters/config"

# Batch
run_tests.bat -package "websocket/internal/adapters/config"
```

Run a specific test:

```powershell
# PowerShell
.\run_tests.ps1 -package "websocket/internal/adapters/config" -test "TestIsDevMode"

# Batch
run_tests.bat -package "websocket/internal/adapters/config" -test "TestIsDevMode"
```

Generate coverage profile and HTML report:

```powershell
# PowerShell
.\run_tests.ps1 -coverprofile "coverage.out"

# Batch
run_tests.bat -coverprofile "coverage.out"
```

## Troubleshooting

If you encounter test failures, try the following:

1. Run tests with verbose output to get more information about the failures:
   ```powershell
   # PowerShell
   .\run_tests.ps1 -v

   # Batch
   run_tests.bat -v
   ```

2. Run tests for a specific package to isolate the issue:
   ```powershell
   # PowerShell
   .\run_tests.ps1 -package "websocket/internal/adapters/config" -v

   # Batch
   run_tests.bat -package "websocket/internal/adapters/config" -v
   ```

3. Fix any issues in the test files and run the tests again.

4. If you're having issues with the coverage profile, try running the commands separately:
   ```bash
   # Run the test with coverage profile
   go test -coverprofile=coverage.out ./...

   # Generate the HTML report
   go tool cover -html=coverage.out -o coverage.html
   ```
