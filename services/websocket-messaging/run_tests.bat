@echo off
setlocal enabledelayedexpansion

REM Run all tests in the project
REM Usage: run_tests.bat [options]
REM Options:
REM   -v: Verbose output
REM   -cover: Show coverage
REM   -coverprofile <file>: Generate coverage profile
REM   -package <package>: Run tests for a specific package
REM   -test <test>: Run a specific test

set cmd=go test
set verbose=0
set cover=0
set coverprofile=
set package=
set test=

:parse_args
if "%~1"=="" goto run_tests
if "%~1"=="-v" (
    set verbose=1
    shift
    goto parse_args
)
if "%~1"=="-cover" (
    set cover=1
    shift
    goto parse_args
)
if "%~1"=="-coverprofile" (
    set coverprofile=%~2
    shift
    shift
    goto parse_args
)
if "%~1"=="-package" (
    set package=%~2
    shift
    shift
    goto parse_args
)
if "%~1"=="-test" (
    set test=%~2
    shift
    shift
    goto parse_args
)
shift
goto parse_args

:run_tests
if %verbose%==1 (
    set cmd=%cmd% -v
)
if %cover%==1 (
    set cmd=%cmd% -cover
)
if not "%coverprofile%"=="" (
    set cmd=%cmd% -coverprofile=%coverprofile%
)
if not "%package%"=="" (
    set cmd=%cmd% %package%
) else (
    set cmd=%cmd% ./...
)
if not "%test%"=="" (
    set cmd=%cmd% -run %test%
)

echo Running: %cmd%
%cmd%

if not "%coverprofile%"=="" (
    echo Generating HTML coverage report...
    go tool cover -html=%coverprofile% -o coverage.html
    echo Coverage report generated: coverage.html
)

endlocal