#!/bin/bash
set -e

GITHUB_ORG="aes128-dev"
GITHUB_REPO="aes128-cli"
SINGBOX_VERSION="1.11.15"
INSTALL_PATH="/usr/local/bin/aes128-cli"
CORE_INSTALL_DIR="/usr/lib/aes128-cli"

main() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "Error: This script must be run as root. Please run with sudo." >&2
        exit 1
    fi
    
    check_dependencies
    
    if [ -f "$INSTALL_PATH" ]; then
        echo "--> Existing aes128-cli installation found. Updating..."
        update_app
    else
        echo "--> No existing installation found. Starting new installation..."
        install_app
    fi
}

check_dependencies() {
    echo "--> Checking for required tools..."
    local missing_tools=0
    
    if ! command -v curl >/dev/null 2>&1; then
        echo "Error: 'curl' is not installed." >&2
        missing_tools=1
    fi
    if ! command -v tar >/dev/null 2>&1; then
        echo "Error: 'tar' is not installed." >&2
        missing_tools=1
    fi
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    if [ "$OS" = "linux" ] && ! command -v setcap >/dev/null 2>&1; then
        echo "Error: 'setcap' is not installed." >&2
        missing_tools=1
    fi
    
    if [ "$missing_tools" -eq 1 ]; then
        echo ""
        echo "Please install the missing tools and run this script again."
        if [ -f /etc/debian_version ]; then
            echo "On Debian/Ubuntu, run: sudo apt update && sudo apt install curl libcap2-bin"
        elif [ -f /etc/redhat-release ]; then
            echo "On Fedora/CentOS, run: sudo dnf install curl libcap"
        elif [ -f /etc/arch-release ]; then
            echo "On Arch Linux, run: sudo pacman -S curl libcap"
        fi
        exit 1
    fi
}

install_app() {
    ARCH=$(uname -m)
    echo "--> Detected OS: $OS"
    echo "--> Detected Arch: $ARCH"
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="armv7" ;;
        i386 | i686) ARCH="386" ;;
    esac

    install_binaries "$OS" "$ARCH"
    setup_service "$OS"
    
    echo ""
    echo "--> aes128-cli has been successfully installed!"
    echo "You can now use the 'aes128-cli' command globally."
    if [ "$OS" = "linux" ]; then
        echo "To manage the VPN service, use:"
        echo "  sudo systemctl start aes128-cli"
        echo "  sudo systemctl stop aes128-cli"
        echo "  sudo systemctl status aes128-cli"
    fi
}

update_app() {
    ARCH=$(uname -m)
    install_binaries "$OS" "$ARCH"

    if [ "$OS" = "linux" ] && systemctl is-active --quiet aes128-cli; then
        echo "--> Restarting aes128-cli service..."
        systemctl restart aes128-cli
    fi

    echo ""
    echo "--> aes128-cli has been successfully updated!"
}

install_binaries() {
    local OS="$1"
    local ARCH="$2"
    
    echo "--> Downloading and installing binaries for $OS/$ARCH..."
    
    CLI_URL="https://github.com/$GITHUB_ORG/$GITHUB_REPO/releases/latest/download/aes128-cli-$OS-$ARCH"
    SINGBOX_URL="https://github.com/SagerNet/sing-box/releases/download/v$SINGBOX_VERSION/sing-box-$SINGBOX_VERSION-$OS-$ARCH.tar.gz"

    echo "--> Downloading latest aes128-cli..."
    curl -sSL "$CLI_URL" -o "$INSTALL_PATH"
    chmod +x "$INSTALL_PATH"
    
    if [ "$OS" = "linux" ]; then
        echo "--> Setting capabilities for ping..."
        setcap cap_net_raw+ep "$INSTALL_PATH"
    fi

    mkdir -p "$CORE_INSTALL_DIR"
    
    echo "--> Downloading sing-box core..."
    curl -sSL "$SINGBOX_URL" | tar -xz -C "$CORE_INSTALL_DIR"
    
    find "$CORE_INSTALL_DIR" -name "sing-box" -type f -exec mv {} "$CORE_INSTALL_DIR/core" \;
    chmod +x "$CORE_INSTALL_DIR/core"
    
    find "$CORE_INSTALL_DIR" -type d -empty -delete
}

setup_service() {
    local OS="$1"
    echo "--> Setting up system service..."
    
    SUDO_USER=${SUDO_USER:-$(who am i | awk '{print $1}')}
    if [ -z "$SUDO_USER" ]; then
        echo "Error: Could not determine the user who ran sudo." >&2
        exit 1
    fi

    if [ "$OS" = "linux" ] && [ -d /etc/systemd/system ]; then
        echo "--> Creating systemd service file..."
        cat > /etc/systemd/system/aes128-cli.service <<EOF
[Unit]
Description=AES128 CLI VPN Service
After=network-online.target
Wants=network-online.target
[Service]
Type=simple
User=$SUDO_USER
ExecStart=/usr/local/bin/aes128-cli connect
ExecStop=/usr/local/bin/aes128-cli disconnect
Restart=on-failure
RestartSec=5
[Install]
WantedBy=multi-user.target
EOF
        systemctl daemon-reload
        echo "--> To enable the service to start on boot, run: sudo systemctl enable aes128-cli"
    fi
}

main "$@"