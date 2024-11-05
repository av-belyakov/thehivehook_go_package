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

#### Подключение к NATS

GO_HIVEHOOK_NHOST
GO_HIVEHOOK_NPORT
GO_HIVEHOOK_NSUBSCRIBERS - данный параметр должен содержать строку вида:
<наименование события>:<наименование абонента 1>,<наименование абонента 2>;<наименование события>:<наименование абонента 1>
например:
caseupdate:gcm,rcmmsk,rcmnvs;alertupdate:gcm

#### Подключение к TheHive

GO_HIVEHOOK_THHOST
GO_HIVEHOOK_THPORT
GO_HIVEHOOK_THAPIKEY - ЭТО ОБЯЗАТЕЛЬНЫЙ ПАРАМЕТР!!!
Он задается ТОЛЬКО через переменную окружения. В конфигурационном
файле этого параметра нет.

#### Настройки WebHook сервера

GO_HIVEHOOK_WEBHNAME //наименование сервера (gcm, rcmnvs и т.д.)
GO_HIVEHOOK_WEBHHOST
GO_HIVEHOOK_WEBHPORT
GO_HIVEHOOK_WEBHTTLTMPINFO //время жизни временной информации, в секундах от 10 до 86400

#### Настройки достапа к SQLite

GO_HIVEHOOK_SLPATHDB //путь к файлу базы данных SQLite, обязательный параметр так как база SQLite выступает в роли кеширующего буфера для хранения команд от сторонних систем к TheHive

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов.

## Примеры команд передаваемые TheHiveHook_go_Package

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
