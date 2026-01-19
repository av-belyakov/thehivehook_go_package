# Thehivehook_go_package

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

Пакет 'thehivehook_go_package' раелизует сервис который является посредником между TheHive и NATS, и применяется для передачи событий типа 'case' и 'alert' от TheHive через брокер сообщений NATS сторонним сервисам.

## Конфигурационные настройки

Конфигурационные параметры для сервиса, могут быть заданы как через конфигурационный файл, так и методом установки переменных окружения. Однако, все пароли и ключевые токены, используемые для авторизации, задаются ТОЛЬКО через переменные окружения.

### Развёртывание пакета через CI/CD

У каждого центра мониторинга на _gitlab.cloud.gcm_ есть репозиторий **\_thehivehook_go**

```
_thehivehook_go/
├── config
│   ├── config_dev.yml
│   ├── config_prod.yml
│   └── config.yml
├── docker-compose.yml
└── README.md
```

Необходимо выполнить следующие действия:

1. Поверить и при необходимости изменить основной файл настроек _config_prod.yml_.

```yaml
COMMONINFO:
  file_name: config_prod
NATS:
  prefix: test
  host: nats.cloud.gcm
  port: 4222
  cache_ttl: 3600
  subscriptions:
    sender_case: object.casetype
    sender_alert: object.alerttype
    listener_command: object.commandstype
THEHIVE:
  host: thehive.cloud.gcm #ip или домен своего TheHive
  port: 9000
  cache_ttl: 43200
WEBHOOKSERVER:
  name: rcmnvs #наименование webhookserver (у каждого своё)
  host: git-runner.rcm #ip или домен ВМ на которой будет хостится приложение
  port: 5000
  ttl_tmp_info: 10
DATABASEWRITELOG:
  host: datahook.cloud.gcm #домен сервера БД для логов (не менять, домен ГЦМ)
  port: 9200
  namedb:
  storage_name_db: thehivehook_go_package
  user: log_writer
```

Важно проверить актуальность параметров **THEHIVE.host**, **THEHIVE.port**, параметр **WEBHOOKSERVER.host** может быть пустым. Параметр **WEBHOOKSERVER.name** задаётся центром мониторинга. Использовать сокращённое, общепризнанное в ведомстве имя. Остальные параметры файла менять не надо.
Файл _config_dev.yml_ используется для тестов в рамках ГЦМ, править не надо. Он применяется только когда переменна окружения **GO_HIVEHOOK_MAIN=development**.
Файл _config.yml_ содержит настройки записи и передачи логов, а так же взаимодействия с Zabbix. Особого смысла править нет, если конечно нет желания отключить логирования некоторых событий и изменить место или размер их хранения.

2. Проверить и при необходимости актулизировать значение **HIVEHOOK_THAPIKEY** содержащее токен доступа к API TheHive. Находится в _Setting->CI/CD->Variables_, если значения там нет - создать. Значение **HIVEHOOK_DBWLOGPASSWD** править не надо.

3. Создать и развернуть gitlab-runner, http://gitlab.cloud.gcm/aguslikov/help в помощь.

#### Типы конфигурационных файлов:

- _config.yaml_ общий конфигурационный файл;
- _config_dev.yaml_ конфигурационный файл используемый для тестов при отладке;
- _config_test.yaml_ конфигурационный файл используемый для тестов при разработке;
- _config_prod.yaml_ конфигурационный файл применяемый в продуктовом режиме.

Основная переменная окружения для данного приложения - **GO_HIVEHOOK_MAIN**. На основании значения этой переменной принимается решение какой из конфигурационных файлов _config_dev.yaml_ или _config_prod.yaml_ использовать. При **GO_HIVEHOOK_MAIN=development** будет использоваться _config_dev.yaml_, при **GO_HIVEHOOK_MAIN=test** будет использоваться _config_test.yaml_, во всех остальных случаях, в том числе и при отсутствии переменной окружения **GO_HIVEHOOK_MAIN** будет использоваться конфигурационный файл _config_prod.yaml_. Перечень переменных окружения которые можно использовать для настройки приложения:

#### Переменная окружения отвечающая за тип запуска приложения "test", "development" или "production"

- **GO_HIVEHOOK_MAIN**

#### Переменные окружения отвечающие за подключение к NATS

- **GO_HIVEHOOK_NHOST**
- **GO_HIVEHOOK_NPORT**
- **GO_HIVEHOOK_NCACHETTL** - данный параметр должен содержать время жизни записи
  кэша, по истечение которого запись автоматически удаляется, значение задается
  в секундах в диапазоне от 10 до 86400 секунд
- **GO_HIVEHOOK_NSUBSENDERCASE** - канал для отправки в него информации по case
- **GO_HIVEHOOK_NSUBSENDERALERT** - канал для отправки в него информации по alert
- **GO_HIVEHOOK_NSUBLISTENERCOMMAND** - канал для приема команд которые нужно выполнить на TheHive

#### Переменные окружения отвечающие за подключение к TheHive

- **GO_HIVEHOOK_THHOST**
- **GO_HIVEHOOK_THPORT**
- **GO_HIVEHOOK_THCACHETTL** - данный параметр должен содержать время жизни записи
  кэша, по истечение которого запись автоматически удаляется, значение задается
  в секундах в диапазоне от 10 до 86400 секунд
- **GO_HIVEHOOK_THAPIKEY** - ЭТО ОБЯЗАТЕЛЬНЫЙ ПАРАМЕТР!!!
  Он задается ТОЛЬКО через переменную окружения. В конфигурационном
  файле этого параметра нет.

#### Переменные окружения отвечающие за настройки WebHook сервера

- **GO_HIVEHOOK_WEBHNAME** //наименование сервера (gcm, rcmnvs и т.д.)
- **GO_HIVEHOOK_WEBHHOST**
- **GO_HIVEHOOK_WEBHPORT**
- **GO_HIVEHOOK_WEBHTTLTMPINFO** //время жизни временной информации, в секундах от 10 до 86400

#### Переменные окружения отвечающие за настройки доступа к БД в которую будут записыватся логи

- **GO_HIVEHOOK_DBWLOGHOST** // доменное имя или ip БД
- **GO_HIVEHOOK_DBWLOGPORT** // порт БД
- **GO_HIVEHOOK_DBWLOGNAME** // наименование БД (при необходимости)
- **GO_HIVEHOOK_DBWLOGSTORAGENAME** // наименование объекта хранения логов (таблица, документ, индекс и т.д., зависит от типа БД)
- **GO_HIVEHOOK_DBWLOGUSER** // пользователь БД
- **GO_HIVEHOOK_DBWLOGPASSWD** // пароль для доступа к БД

Настройки логирования данных в БД не являются обязательными и необходимы только если пользователь приложения желает хранить логи в базе данных

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов.

## Профилирование приложения

Профилирование приложения возможно только в двух режимах "test" или "development". Для того что бы получить доступ к профилировщику нужно выполнить в браузере или GET запросы с помощью wget или curl следующее:

```bash
http://ip:port/debug/pprof/
```

Использование инструмента 'go tool pprof'.

```bash
go tool pprof http://ip:port/debug/pprof/... (далее возможны вариации heap, allocs, goroutine и т.д.)
```

- **goroutine** — следы всех текущих горутин;
- **heap** — выборка выделений памяти живых объектов;
- **allocs** — выборка всех прошлых выделений памяти;
- **threadcreate** — следы стека, которые привели к созданию новых потоков в операционной системе;
- **block** — следы стека, которые привели к блокировке примитивов синхронизации;
- **mutex** — следы стека держателей конфликтующих мьютексов.

## Примеры команд передаваемые 'thehivehook_go_package'

Все команды для 'thehivehook_go_package' представляют собой JSON объекты передаваемые в бинарном виде. Структура и значение команд обрабатываемых 'thehivehook_go_package':

```json
{
  "service": "наименование сервиса", //обязательный параметр
  "command": "команда",
  "for_regional_object": "имя регионального thehivehook",
  "root_id": "основной id, как правило это rootId case или alert", //обязательный параметр только для некоторых действий выполняемых с конкретным кейсом или алертом
  "case_id": "id кейса, если есть",
  "username": "имя пользователя", //необходим если нужно указать пользователя выполнившего действие
  "field_name": "некое ключевое поле",
  "value": "устанавливаемое значение",
  "byte_data": "набор данных в бинарном виде"
}
```

#### Перечень видов обрабатываемых команд:

- **add_case_tag**
- **add_case_task**
- **set_case_custom_field**

Пример команды для добавления тега:

```json
{
  "service": "MISP",
  "command": "add_case_tag",
  "for_regional_object": "gcm",
  "root_id": "~74395656",
  "case_id": "13435",
  "value": "Webhook: send=\"MISP\""
}
```

Пример команды для добавления задачи:

```json
{
  "service": "MISP",
  "command": "add_case_task",
  "for_regional_object": "имя регионального thehivehook",
  "root_id": "~74395656",
  "username": "architector@33c.rcm",
  "field_name": "Developers",
  "value": "handling request"
}
```

Пример команды для добавления поля custom field:

```json
{
  "service": "MISP",
  "command": "set_case_custom_field",
  "for_regional_object": "имя регионального thehivehook",
  "root_id": "~74395656",
  "field_name": "misp-event-id.string",
  "value": "3221"
}
```

По полю **for_regional_object** JSON объекта модуль определяет ему ли предназначена полученная команда. Если имя в поле **for_regional_object** не соответсвует имени **WEBHOOKSERVER.name** в конфигурационных файлах _config_test_, _config_dev_ или _config_prod_, использование одного из этих файлов зависит от значения переменной окружения **GO_HIVEHOOK_MAIN**, то модуль записывает в лог файл _error.log_ соответствующее сообщение. Поле **error** JSON ответа будет содержать сообщение об ошибке, а поле **status_code** код ответа '400'.

#### Структура ответа на любую из переданных команд

Пример ответа:

```json
{
  "id": "",
  "source": "gcm",
  "error": "no error",
  "command": "",
  "status_code": 0,
  "data": "дополнительные данные, возможны в любом типе"
}
```

Ответ отправляется на любую команду, в том числе и на неверно сформированный JSON запрос.

## Настройка 'endpoints' для TheHive

Добавить в конфигурационный файл TheHive (thehive/conf/application.conf) в параметр notification.webhook.endpoints значение с новым 'endpoint', по аналогии.

Далее нужно выполнить перезагрузку TheHive что бы применился поправленый конфиг.

Следом выполняем:

```bash
curl -XPUT -H "Authorization: Bearer <ApiKey>" -H 'Content-type: application/json' <url*или*ip*и*сетевой*порт>/api/config/organisation/notification -d
  '{
    "value": [
      {
        "delegate": false,
        "trigger": {"name": "AnyEvent"},
        "notifier": {"name": "webhook", "endpoint": "hivehook"}
      }
    ]
  }'
```

Что бы посмотреть все доступные endpoints на хайве нужно выполнить:

```bash
curl  -H "Authorization: Bearer <api_key_local_admin>" http://**IP ВАШЕГО Thehive**:**Порт ВАШЕГО Thehive**/api/config/notification.webhook.endpoints
```

однако для этого нужны привелегии пользователя admin.
