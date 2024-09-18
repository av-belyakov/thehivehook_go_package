# README

Приложение "TheHiveHook_go_package" преднозначено для получения данных от TheHive, обогащения этих данных
и передачи их через NATS любым другим заинтересованным системам

## Быстрый старт

## Структура поекта

## Конфигурационные настройки

Конфигурационные параметры для сервиса могут быть заданы как через конфигурационный файл так и методом установки переменных окружения. Однако, все пароли и
ключевые токены, используемые для авторизации, задаются ТОЛЬКО через переменные окружения.

Типы конфигурационных файлов:

- config.yaml общий конфигурационный файл
- config_dev.yaml конфигурационный файл используемый для тестов при разработке
- config_prod.yaml конфигурационный файл применяемый в продуктовом режиме

Основная переменная окружения для данного приложения - GO_HIVEHOOK_MAIN. На основании значения этой переменной принимается решение какой из конфигурационных файлов config_dev.yaml или config_prod.yaml использовать. При GO_HIVEHOOK_MAIN=development будет использоваться config_dev.yaml, во всех остальных случаях, в том числе и при отсутствии переменной окружения GO_HIVEHOOK_MAIN будет использоваться конфигурационный файл config_prod.yaml. Перечень переменных окружения которые можно использовать для настройки приложения:

//Переменная окружения отвечающая за тип запуска приложения "development" или "production"
GO_HIVEHOOK_MAIN

//Подключение к NATS
GO_HIVEHOOK_NHOST
GO_HIVEHOOK_NPORT
GO_HIVEHOOK_NSUBSCRIBERS - данный параметр должен содержать строку вида:
<наименование события>:<наименование абонента 1>,<наименование абонента 2>;<наименование события>:<наименование абонента 1>
например:
caseupdate:gcm,rcmmsk,rcmnvs;alertupdate:gcm

//Подключение к TheHive
GO_HIVEHOOK_NAME
GO_HIVEHOOK_URL

GO_HIVEHOOK_APIKEY - ЭТО ОБЯЗАТЕЛЬНЫЙ ПАРАМЕТР!!!
Он задается ТОЛЬКО через переменную окружения. В конфигурационном
файле этого параметра нет.

//Подключение к СУБД Elasticsearch
GO_HIVEHOOK_ESHOST
GO_HIVEHOOK_ESPORT
GO_HIVEHOOK_ESUSER

GO_HIVEHOOK_ESPASSWD - ЭТО ОБЯЗАТЕЛЬНЫЙ ПАРАМЕТР!!! Он задается ТОЛЬКО через переменную окружения. В конфигурационном
файле этого параметра нет.

GO_HIVEHOOK_ESPREFIX
GO_HIVEHOOK_ESINDEX

//Настройки основного API сервера
GO_HIVEHOOK_HHOST
GO_HIVEHOOK_HPORT

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов.
