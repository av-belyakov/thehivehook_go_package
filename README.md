# TheHiveHook_Go_Package

TheHiveHook_go_package является посредником между TheHive и NATS и применяется для передачи кейсов и алертов TheHive, через брокер сообщений NATS, любым другим заинтересованным системам

## Конфигурационные настройки

Конфигурационные параметры для сервиса могут быть заданы как через конфигурационный файл так и методом установки переменных окружения. Однако, все пароли и
ключевые токены, используемые для авторизации, задаются ТОЛЬКО через переменные окружения.

#### Типы конфигурационных файлов:

- config.yaml общий конфигурационный файл
- config_dev.yaml конфигурационный файл используемый для тестов при разработке
- config_prod.yaml конфигурационный файл применяемый в продуктовом режиме

Основная переменная окружения для данного приложения - GO_HIVEHOOK_MAIN. На основании значения этой переменной принимается решение какой из конфигурационных файлов config_dev.yaml или config_prod.yaml использовать. При GO_HIVEHOOK_MAIN=development будет использоваться config_dev.yaml, во всех остальных случаях, в том числе и при отсутствии переменной окружения GO_HIVEHOOK_MAIN будет использоваться конфигурационный файл config_prod.yaml. Перечень переменных окружения которые можно использовать для настройки приложения:

#### Переменная окружения отвечающая за тип запуска приложения "development" или "production"

GO_HIVEHOOK_MAIN

#### Переменные окружения отвечающие за подключение к NATS

GO_HIVEHOOK_NPREFIX
GO_HIVEHOOK_NHOST
GO_HIVEHOOK_NPORT
GO_HIVEHOOK_NCACHETTL - данный параметр должен содержать время жизни записи
кэша, по истечение которого запись автоматически удаляется, значение задается
в секундах в диапазоне от 10 до 86400 секунд
GO_HIVEHOOK_NSUBSENDERCASE - канал для отправки в него информации по case
GO_HIVEHOOK_NSUBSENDERALERT - канал для отправки в него информации по alert
GO_HIVEHOOK_NSUBLISTENERCOMMAND - канал для приема команд которые нужно выполнить на TheHive

#### Переменные окружения отвечающие за подключение к TheHive

GO_HIVEHOOK_THHOST
GO_HIVEHOOK_THPORT
GO_HIVEHOOK_THCACHETTL - данный параметр должен содержать время жизни записи
кэша, по истечение которого запись автоматически удаляется, значение задается
в секундах в диапазоне от 10 до 86400 секунд
GO_HIVEHOOK_THAPIKEY - ЭТО ОБЯЗАТЕЛЬНЫЙ ПАРАМЕТР!!!
Он задается ТОЛЬКО через переменную окружения. В конфигурационном
файле этого параметра нет.

#### Переменные окружения отвечающие за настройки WebHook сервера

GO_HIVEHOOK_WEBHNAME //наименование сервера (gcm, rcmnvs и т.д.)
GO_HIVEHOOK_WEBHHOST
GO_HIVEHOOK_WEBHPORT
GO_HIVEHOOK_WEBHTTLTMPINFO //время жизни временной информации, в секундах от 10 до 86400

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов.

## Примеры команд передаваемые TheHiveHook_Go_Package

Все команды для TheHiveHook_Go_Package представляют собой JSON объекты передаваемые
в бинарном виде. Структура и значение команд обрабатываемых TheHiveHook_Go_Package:

```
{
  "service": "<наименование сервиса>" //обязательный параметр
  "commands": [
    {
      "command": <команда>
      "root_id": "<основной id, как правило это rootId case или alert>" //обязательный параметр только для некоторых действий выполняемых с конкретным кейсом или алертом
      "case_id": "<id кейса, если есть>"
      "username": <имя пользователя> //необходим если нужно указать пользователя выполнившего действие
      "field_name": <некое ключевое поле>
      "value": <устанавливаемое значение>
      "byteData": <набор данных в бинарном виде>
    }
  ]
}
```

##### Перечень видов обрабатываемых команд:

- "add_tag"
- "add_task"
- "set_custom_field"

Пример команды для добавления тега:

```
{
  "service": "MISP",
  "commands": [
    {
      "command": "add_tag",
      "root_id": "~74395656",
      "case_id": "13435",
      "value": "Webhook: send=\"MISP\""
    }
  ]
}
```

Пример команды для добавления задачи:

```
{
  "service": "MISP",
  "commands": [
    {
      "command": "add_task",
      "root_id": "~74395656",
      "username": "architector@33c.rcm",
      "field_name": "Developers",
      "value": "handling request"
    }
  ]
}
```

Пример команды для добавления поля custom field:

```
{
  "service": "MISP",
  "commands": [
    {
      "command": "set_custom_field",
      "root_id": "~74395656",
      "field_name": "misp-event-id.string",
      "value": "3221"
    }
  ]
}
```

Пример набора команд для объектов типа Case, где команда

- "addtag" добавляет тег если такого тега еще нет
- "setcustomfield" устанавливает настраиваемые поля
- "addtask" добавляет задачи
  `{
  "success": true,
  "service": "SERVICE_NAME",
  "error": "ERROR_IF_EXISTS",
  "commands": [
    {
      "command": "addtag",
      "string": "Webhook: send=\"MISP\""
    },
    {
      "command": "setcustomfield",
      "name": "misp-event-id.string",
      "string": "123"
    },
    {
      "command": "addtask",
      "name": "Developers",
      "string": "not added",
      "username": "architector@33c.rcm"
    }
  ]
}`

##

## Настройка TheHive

curl -XPUT -u <имя*пользователя>:'пароль_org-admin' -H 'Content-type: application/json' <url*или*ip*и*сетевой*порт>/api/config/organisation/notification -d '
{
"value": [
{
"delegate": false,
"trigger": { "name": "AnyEvent"},
"notifier": { "name": "webhook", "endpoint": "hivehook" }
}
]
}'
