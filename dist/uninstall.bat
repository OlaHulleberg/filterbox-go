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

:: Stop and uninstall the daemon
echo "Stopping the daemon..."
cmd /c "filterbox-daemon.exe stop"
echo "Uninstalling the daemon..."
cmd /c "filterbox-daemon.exe uninstall"

:: Now remove the files
echo "Removing the files..."
cmd /c "del filterbox-daemon.log"
cmd /c "del filterbox-daemon.exe"
cmd /c "del filterbox.exe"

echo "FilterBox uninstalled successfully."

:: Remove the installation directory and the script itself
cd ..
rmdir /s /q "FilterBox"

:: Exit the script
exit /b