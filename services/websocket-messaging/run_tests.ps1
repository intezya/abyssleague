# Run all tests in the project
# Usage: .\run_tests.ps1 [options]
# Options:
#   -v: Verbose output
#   -cover: Show coverage
#   -coverprofile: Generate coverage profile
#   -package: Run tests for a specific package (e.g. "websocket/internal/adapters/config")
#   -test: Run a specific test (e.g. "TestIsDevMode")

param (
    [switch]$v,
    [switch]$cover,
    [string]$coverprofile = "",
    [string]$package = "",
    [string]$test = ""
)

# Build the command
$cmd = "go test"

if ($v) {
    $cmd += " -v"
}

if ($cover) {
    $cmd += " -cover"
}

if ($coverprofile -ne "") {
    $cmd += " -coverprofile=$coverprofile"
}

if ($package -ne "") {
    $cmd += " $package"
} else {
    $cmd += " ./..."
}

if ($test -ne "") {
    $cmd += " -run $test"
}

# Run the command
Write-Host "Running: $cmd"
Invoke-Expression $cmd

# If coverage profile was generated, show HTML report
if ($coverprofile -ne "") {
    Write-Host "Generating HTML coverage report..."
    Invoke-Expression "go tool cover -html=$coverprofile -o coverage.html"
    Write-Host "Coverage report generated: coverage.html"
}