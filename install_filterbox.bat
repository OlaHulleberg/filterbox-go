@echo off
:: Batch script to build and install FilterBox on Windows using Makefile

:: Check for administrative privileges (if not, restart the script with them)
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting administrative privileges...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

:: Set the current directory to the script directory
cd %~dp0

:: Run make build
echo Building project...
make all
if %errorLevel% neq 0 (
    echo Failed to build project.
    pause
    exit /b %errorLevel%
)

:: Run make install
echo Installing project...
make install
if %errorLevel% neq 0 (
    echo Failed to install project.
    pause
    exit /b %errorLevel%
)

echo FilterBox installed successfully.
pause