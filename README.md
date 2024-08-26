# README

Приложение "TheHiveHook_go_package" преднозначено для получения данных от TheHive, обогащения этих данных
и передачи их через NATS любым другим заинтересованным системам

## Быстрый старт

## Структура поекта

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|                           |                                                                                                                                                                                        |
| ------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`cmd/api`**             | Your application-specific code (handlers, routing, middleware, helpers) for dealing with HTTP requests and responses.                                                                  |
| `↳ cmd/api/errors.go`     | Contains helpers for managing and responding to error conditions.                                                                                                                      |
| `↳ cmd/api/handlers.go`   | Contains your application HTTP handlers.                                                                                                                                               |
| `↳ cmd/api/helpers.go`    | Contains helper functions for common tasks.                                                                                                                                            |
| `↳ cmd/api/main.go`       | The entry point for the application. Responsible for parsing configuration settings initializing dependencies and running the server. Start here when you're looking through the code. |
| `↳ cmd/api/middleware.go` | Contains your application middleware.                                                                                                                                                  |
| `↳ cmd/api/routes.go`     | Contains your application route mappings.                                                                                                                                              |
| `↳ cmd/api/server.go`     | Contains a helper functions for starting and gracefully shutting down the server.                                                                                                      |

|                         |                                                           |
| ----------------------- | --------------------------------------------------------- |
| **`internal`**          | Contains various helper packages used by the application. |
| `↳ internal/request/`   | Contains helper functions for decoding JSON requests.     |
| `↳ internal/response/`  | Contains helper functions for sending JSON responses.     |
| `↳ internal/validator/` | Contains validation helpers.                              |
| `↳ internal/version/`   | Contains the application version number definition.       |

## Конфигурационные настройки

Configuration settings are managed via command-line flags in `main.go`.

You can try this out by using the `--http-port` flag to configure the network port that the server is listening on:

```
$ go run ./cmd/api --http-port=9999
```

Feel free to adapt the `run()` function to parse additional command-line flags and store their values in the `config` struct. For example, to add a configuration setting to enable a 'debug mode' in your application you could do this:

```
type config struct {
    httpPort  int
    debug     bool
}

...

func run() {
    var cfg config

    flag.IntVar(&cfg.httpPort, "http-port", 4444, "port to listen on for HTTP requests")
    flag.BoolVar(&cfg.debug, "debug", false, "enable debug mode")

    flag.Parse()

    ...
}
```
