package sqliteapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// New инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func New(ctx context.Context, path string, logging *logginghandler.LoggingChan) (chan<- ChanApiSqlite, error) {
	chanListene := make(chan ChanApiSqlite)

	if path == "" {
		return chanListene, errors.New("the path to the database file should not be empty")
	}

	settings := apiSqliteSettings{
		path:   path,
		logger: logging,
	}

	sqldb, err := sql.Open("sqlite3", settings.path)
	if err != nil {
		return chanListene, err
	}

	if sqldb.Ping() != nil {
		return chanListene, err
	}

	settings.db = sqldb

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-chanListene:
				settings.route(routeSettings{
					command:      msg.Command,
					taskId:       msg.TaskID,
					service:      msg.Service,
					data:         msg.Data,
					chanResponse: msg.ChanResponse,
				})
			}
		}
	}()

	return chanListene, nil
}

// Route маршрутизатор обработки запросов
func (apisqlite *apiSqliteSettings) route(settings routeSettings) {
	if settings.taskId == "" || settings.command == "" {
		_, f, l, _ := runtime.Caller(0)
		apisqlite.logger.Send("error", fmt.Sprintf(" 'the sql query cannot be processed, the command and the task ID must not be empty' %s:%d", f, l-1))

		return
	}

	switch settings.command {
	case "insert section tags":
		go apisqlite.handlerSectionInsertTags(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "insert section creater":
		go apisqlite.handlerSectionInsertCreater(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "select section tags":
		go apisqlite.handlerSectionSelectTags(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "select section creater":
		go apisqlite.handlerSectionSelectCreater(settings.taskId, settings.service, settings.data, settings.chanResponse)
	}
}

// handlerSectionInsertTags обработчик секции добавления информации о тегах
func (apisqlite *apiSqliteSettings) handlerSectionInsertTags(taskId, service string, data []byte, chanResponse chan<- ChanOutputApiSqlite) {
	response := ChanOutputApiSqlite{TaskID: taskId}

	stmt := apisqlite.prepareTableExecutedCommands()

	//порядок: id, service, binary_data
	result, err := stmt.Exec(taskId, service, data)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		apisqlite.logger.Send("error", fmt.Sprintf(" '%w' %s:%d", err, f, l-1))

		chanResponse <- response

		return
	}

	fmt.Println("func 'handlerSectionInsertTags' RESPONSE:", result)

	//
	// Вообще лучше сделать что бы тип chanResponse
	// соответствовал определенному интерфейсу и далее описать этот
	// интерфейс в commoninterfaces.interfaces.go
	// так будет полее гибко
	//

	response.Status = true
	chanResponse <- response
}

// handlerSectionInsertCreater обработчик секции добавления информации о создаваемых, новых данных
func (apisqlite *apiSqliteSettings) handlerSectionInsertCreater(taskId, service string, data []byte, chanResponse chan<- ChanOutputApiSqlite) {

}

// handlerSectionSelectTags обработчик секции запроса информации о тегах
func (apisqlite *apiSqliteSettings) handlerSectionSelectTags(taskId, service string, data []byte, chanResponse chan<- ChanOutputApiSqlite) {

}

// handlerSectionSelectCreater обработчик секции запроса информации о создаваемых, новых данных
func (apisqlite *apiSqliteSettings) handlerSectionSelectCreater(taskId, service string, data []byte, chanResponse chan<- ChanOutputApiSqlite) {

}

func (apisqlite *apiSqliteSettings) prepareTableExecutedCommands() *sql.Stmt {
	stmt, err := apisqlite.db.Prepare("INSERT INTO table_executed_commands (id, service, command, name, description) values(?,?,?,?,?)")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		apisqlite.logger.Send("error", fmt.Sprintf(" '%w' %s:%d", err, f, l-1))

		return nil
	}

	return stmt
}
