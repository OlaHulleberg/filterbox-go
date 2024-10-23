.PHONY: build-ui build-daemon install clean create-install-dir

OS := $(shell uname -s 2>/dev/null || echo Windows_NT)

ifeq ($(OS), Linux)
    BINARY_DIR=./bin
    DAEMON_BINARY=$(BINARY_DIR)/filterbox-daemon
    UI_BINARY=$(BINARY_DIR)/filterbox-ui
    INSTALL_LOCATION=/usr/local/bin
else ifeq ($(OS), Windows_NT)
    BINARY_DIR=bin
    DAEMON_BINARY=$(BINARY_DIR)\filterbox-daemon.exe
    UI_BINARY=$(BINARY_DIR)\filterbox-ui.exe
    INSTALL_LOCATION=C:\Program Files\FilterBox
endif

all: build-daemon build-ui

build-daemon:
	@echo "Building daemon..."
	@go build -o $(DAEMON_BINARY) ./cmd/daemon


build-ui:
	@echo "Building UI..."
	@go build -o $(UI_BINARY) ./cmd/ui

install: create-install-dir
ifeq ($(OS), Linux)
	@echo "Installing on Linux..."
	@sudo cp $(DAEMON_BINARY) ${INSTALL_LOCATION}/filterbox-daemon
	@sudo cp $(UI_BINARY) ${INSTALL_LOCATION}/filterbox
	@echo "Installing systemd service..."
	@cp dist/filterbox-daemon.service ${HOME}/.config/systemd/user/filterbox-daemon.service
	@systemctl --user daemon-reload
	@systemctl --user enable filterbox-daemon
	@systemctl --user start filterbox-daemon
	@echo "FilterBox daemon successfully installed, and started."
else ifeq ($(OS), Windows_NT)
	@echo "Installing on Windows..."
	@copy "$(DAEMON_BINARY)" "$(INSTALL_LOCATION)\filterbox-daemon.exe"
	@copy "$(UI_BINARY)" "$(INSTALL_LOCATION)\filterbox.exe"
	@copy "dist\uninstall.bat" "$(INSTALL_LOCATION)\uninstall.bat"
	@echo Installing Windows service...
	@"$(INSTALL_LOCATION)\filterbox-daemon.exe" install
	@"$(INSTALL_LOCATION)\filterbox-daemon.exe" start
	@echo FilterBox daemon successfully installed and started.
endif

uninstall:
	@echo "Uninstalling"
	@systemctl --user stop filterbox-daemon
	@systemctl --user disable filterbox-daemon
	@sudo rm /usr/local/bin/filterbox-daemon
	@sudo rm /usr/local/bin/filterbox

clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)

create-install-dir:
ifeq ($(OS), Windows_NT)
	@if not exist "$(INSTALL_LOCATION)" mkdir "$(INSTALL_LOCATION)"
endif
