# aes128-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/aes128-dev/aes128-cli)](https://goreportcard.com/report/github.com/aes128-dev/aes128-cli)

A command-line interface (CLI) for the [AES128 VPN](https://aes128.com/en/) service. Allows you to manage your VPN connection directly from the terminal.

## Features

-   Secure authentication and session management.
-   Automatically connect to the fastest available server.
-   Connect to a specific location by its ID or domain.
-   View connection status, including location and uptime.
-   Display a list of available locations with their ping latency.
-   Customize connection protocol (`vless`, `vmess`, `trojan`) and AdBlock settings.
-   Automatic session validation and disconnect on invalid token.

## Installation

Prerequisites: `curl` and `tar`. Run the following command in your terminal:

```bash
curl -sSL https://raw.githubusercontent.com/aes128-dev/aes128-cli/main/install.sh | sudo bash
```

The script will automatically install `aes128-cli` and the required `sing-box` core.

## Usage

**Important:** Commands that manage the system's network state (`connect`, `disconnect`, `status`) require superuser privileges (`sudo`). Commands that handle your account and settings are run without `sudo`.

---

#### `login`
Authenticates the user.
```bash
$ aes128-cli login
```

---
#### `connect [id or domain]`
Connects to the VPN. If no arguments are provided, it connects to the fastest server.
```bash
# Connect to the fastest server
$ sudo aes128-cli connect

# Connect to the server with ID 5
$ sudo aes128-cli connect 5
```

---
#### `disconnect`
Disconnects from the VPN.
```bash
$ sudo aes128-cli disconnect
```

---
#### `status`
Shows the current connection status.
```bash
$ sudo aes128-cli status
```

---
#### `locations`
Shows the list of available locations and their ping.
```bash
$ aes128-cli locations
```

---
#### `account`
Shows information about the current account (username and session name).
```bash
$ aes128-cli account
```

---
#### `sessions`
Manages active sessions.
```bash
# List active sessions
$ aes128-cli sessions list

# Delete a session
$ aes128-cli sessions delete
```

---
#### `settings`
Manages user settings.
```bash
# Show current settings
$ aes128-cli settings get

# Set protocol to trojan
$ aes128-cli settings set protocol trojan

# Enable AdBlock
$ aes128-cli settings set adblock on
```

---
#### `logout`
Logs out and clears all local session data.
```bash
$ aes128-cli logout
```

## Building from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/aes128-dev/aes128-cli.git
    cd aes128-cli
    ```

2.  Run the build script:
    ```bash
    ./build.sh
    ```
    The compiled binaries and checksums will be available in the `release` directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
---

## Русская версия

Командная утилита (CLI) для сервиса [AES128 VPN](https://aes128.com/ru/). Позволяет управлять VPN-соединением напрямую из терминала.

### Возможности

-   Безопасная аутентификация и управление сессиями.
-   Автоматическое подключение к самому быстрому серверу.
-   Подключение к выбранной локации по ID или домену.
-   Просмотр статуса соединения, включая локацию и время работы (uptime).
-   Отображение списка доступных локаций с пингом.
-   Настройка протокола (`vless`, `vmess`, `trojan`) и блокировки рекламы.
-   Автоматическая проверка сессии и отключение при невалидности токена.

### Установка

Для установки требуется `curl` и `tar`. Выполните следующую команду в терминале:

```bash
curl -sSL https://raw.githubusercontent.com/aes128-dev/aes128-cli/main/install.sh | sudo bash
```

Скрипт автоматически установит `aes128-cli` и необходимое ядро `sing-box` в систему.

### Использование

**Важно:** Команды, управляющие системным состоянием сети (`connect`, `disconnect`, `status`), требуют прав суперпользователя (`sudo`). Команды, работающие с аккаунтом и настройками, выполняются без `sudo`.

---

#### `login`
Аутентификация пользователя.
```bash
$ aes128-cli login
```

---
#### `connect [id или домен]`
Подключение к VPN. Без аргументов подключается к самому быстрому серверу.
```bash
# Подключиться к самому быстрому серверу
$ sudo aes128-cli connect

# Подключиться к серверу с ID 5
$ sudo aes128-cli connect 5
```

---
#### `disconnect`
Отключиться от VPN.
```bash
$ sudo aes128-cli disconnect
```

---
#### `status`
Показать текущий статус соединения.
```bash
$ sudo aes128-cli status
```

---
#### `locations`
Показать список доступных локаций и их пинг.
```bash
$ aes128-cli locations
```

---
#### `account`
Показать информацию о текущем аккаунте (имя пользователя и сессии).
```bash
$ aes128-cli account
```

---
#### `sessions`
Управление активными сессиями.
```bash
# Показать список сессий
$ aes128-cli sessions list

# Удалить сессию
$ aes128-cli sessions delete
```

---
#### `settings`
Управление настройками.
```bash
# Показать текущие настройки
$ aes128-cli settings get

# Установить протокол trojan
$ aes128-cli settings set protocol trojan

# Включить блокировку рекламы
$ aes128-cli settings set adblock on
```

---
#### `logout`
Выход из системы и удаление локальных данных сессии.
```bash
$ aes128-cli logout
```

### Сборка из исходного кода

1.  Клонируйте репозиторий:
    ```bash
    git clone https://github.com/aes128-dev/aes128-cli.git
    cd aes128-cli
    ```

2.  Запустите скрипт сборки:
    ```bash
    ./build.sh
    ```
    Готовые бинарные файлы и контрольные суммы появятся в папке `release`.

### Лицензия

Этот проект распространяется под лицензией MIT - подробности смотрите в файле [LICENSE](LICENSE).
