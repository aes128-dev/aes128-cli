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
    
    if [ -f "$INSTALL_PATH" ]; then
        echo "--> Existing aes128-cli installation found. Attempting to update..."
        update_app
    else
        echo "--> No existing aes128-cli installation found. Starting new installation..."
        install_app
    fi
}

install_app() {
    check_dependencies
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
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
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="armv7" ;;
        i386 | i686) ARCH="386" ;;
    esac
    
    install_binaries "$OS" "$ARCH"

    # Перезапускаем сервис, если он был активен
    if systemctl is-active --quiet aes128-cli; then
        echo "--> Restarting aes128-cli service..."
        systemctl restart aes128-cli
    fi

    echo ""
    echo "--> aes128-cli has been successfully updated!"
}

check_dependencies() {
    echo "--> Checking for required tools..."
    if ! command -v curl >/dev/null 2>&1; then
        echo "Error: 'curl' is not installed. Please install it first." >&2
        echo "On Debian/Ubuntu: sudo apt install curl" >&2
        echo "On Fedora/CentOS: sudo dnf install curl" >&2
        exit 1
    fi
    if ! command -v tar >/dev/null 2>&1; then
        echo "Error: 'tar' is not installed. Please install it first." >&2
        exit 1
    fi
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

    mkdir -p "$CORE_INSTALL_DIR"
    
    echo "--> Downloading latest sing-box core..."
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
        
    elif [ "$OS" = "freebsd" ]; then
        echo "--> Creating rc.d service file..."
        cat > /usr/local/etc/rc.d/aes128-cli <<EOF
#!/bin/sh
#
# PROVIDE: aes128-cli
# REQUIRE: LOGIN
# KEYWORD: shutdown
. /etc/rc.subr
name="aes128-cli"
rcvar="aes128_cli_enable"
command="/usr/local/bin/aes128-cli"
start_cmd="\${command} connect"
stop_cmd="\${command} disconnect"
load_rc_config \$name
run_rc_command "\$1"
EOF
        chmod +x /usr/local/etc/rc.d/aes128-cli
        echo "--> To enable the service to start on boot, add aes128_cli_enable=\"YES\" to /etc/rc.conf"
    fi
}

main "$@"