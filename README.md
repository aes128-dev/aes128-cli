# aes128-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/aes128-dev/aes128-cli)](https://goreportcard.com/report/github.com/aes128-dev/aes128-cli)

Командная утилита (CLI) для сервиса AES128 VPN. Позволяет управлять VPN-соединением напрямую из терминала.

## Возможности

-   Безопасная аутентификация и управление сессиями.
-   Автоматическое подключение к самому быстрому серверу.
-   Подключение к выбранной локации по ID или домену.
-   Просмотр статуса соединения, включая локацию и время работы (uptime).
-   Отображение списка доступных локаций с пингом.
-   Настройка протокола (`vless`, `vmess`, `trojan`) и блокировки рекламы.
-   Автоматическая проверка сессии и отключение при невалидности токена.

## Установка

Для установки требуется `curl` и `tar`. Выполните следующую команду в терминале:

```bash
curl -sSL [https://raw.githubusercontent.com/aes128-dev/aes128-cli/main/install.sh](https://raw.githubusercontent.com/aes128-dev/aes128-cli/main/install.sh) | sudo bash
```

Скрипт автоматически установит `aes128-cli` и необходимое ядро `sing-box` в систему.

## Использование

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

## Сборка из исходного кода

1.  Клонируйте репозиторий:
    ```bash
    git clone [https://github.com/aes128-dev/aes128-cli.git](https://github.com/aes128-dev/aes128-cli.git)
    cd aes128-cli
    ```

2.  Запустите скрипт сборки:
    ```bash
    ./build.sh
    ```
    Готовые бинарные файлы и контрольные суммы появятся в папке `release`.

## Лицензия

MIT
