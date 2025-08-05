#!/bin/bash

set -e

main() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "This script must be run as root. Please run with sudo."
        exit 1
    fi

    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    echo "Detected OS: $OS"
    echo "Detected Arch: $ARCH"

    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="armv7" ;;
        i386 | i686) ARCH="386" ;;
    esac

    install_dependencies
    install_binaries
    setup_service
    
    echo ""
    echo "aes128-cli has been successfully installed!"
    echo "You can now use the 'aes128-cli' command globally."
    if [ "$OS" = "linux" ]; then
        echo "To manage the VPN service, use:"
        echo "  sudo systemctl start aes128-cli"
        echo "  sudo systemctl stop aes128-cli"
        echo "  sudo systemctl status aes128-cli"
    fi
}

install_dependencies() {
    echo "Installing required dependencies (curl)..."
    if [ -f /etc/debian_version ]; then
        apt-get update && apt-get install -y curl ca-certificates
    elif [ -f /etc/redhat-release ]; then
        dnf install -y curl
    elif [ -f /etc/arch-release ]; then
        pacman -Sy --noconfirm curl
    elif [ "$OS" = "freebsd" ]; then
        pkg install -y curl ca_root_nss
    else
        echo "Warning: Could not detect package manager. Please ensure curl is installed."
    fi
}

install_binaries() {
    echo "Downloading and installing binaries..."
    
    CLI_URL="https://github.com/aes128-dev/aes128-cli/releases/latest/download/aes128-cli-$OS-$ARCH"
    SINGBOX_URL="https://github.com/SagerNet/sing-box/releases/download/v1.11.15/sing-box-1.11.15-$OS-$ARCH.tar.gz"

    echo "Downloading aes128-cli..."
    curl -sSL "$CLI_URL" -o /usr/local/bin/aes128-cli
    chmod +x /usr/local/bin/aes128-cli

    INSTALL_DIR="/usr/lib/aes128-cli"
    mkdir -p "$INSTALL_DIR"
    
    echo "Downloading sing-box core..."
    curl -sSL "$SINGBOX_URL" | tar -xz -C "$INSTALL_DIR"
    
    find "$INSTALL_DIR" -name "sing-box" -type f -exec mv {} "$INSTALL_DIR/core" \;
    chmod +x "$INSTALL_DIR/core"
    
    find "$INSTALL_DIR" -type d -empty -delete
}

setup_service() {
    echo "Setting up system service..."
    
    SUDO_USER=${SUDO_USER:-$(who am i | awk '{print $1}')}

    if [ "$OS" = "linux" ]; then
        echo "Creating systemd service file..."
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
        echo "To enable the service to start on boot, run: sudo systemctl enable aes128-cli"
        
    elif [ "$OS" = "freebsd" ]; then
        echo "Creating rc.d service file..."
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
        echo "To enable the service to start on boot, add aes128_cli_enable=\"YES\" to /etc/rc.conf"
    fi
}

main "$@"