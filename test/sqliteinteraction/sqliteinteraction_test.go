package sqliteinteraction_test

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testsqliteinteraction", Ordered, func() {
	var (
		db         *sql.DB
		errConnect error
	)

	BeforeAll(func() {
		for k, v := range sql.Drivers() {
			fmt.Printf("%d. %s\n", k+1, v)
		}

		db, errConnect = sql.Open("sqlite3", "../../../../../_sqlite/database/db.sqlite")
	})

	Context("Тест 1. Инициализация и проверка подключения к БД", func() {
		It("При инициализации подключения не должно быть ошибок", func() {
			Expect(errConnect).ShouldNot(HaveOccurred())
		})
		It("При выполнении проверки подключения не должно быть ошибок", func() {
			Expect(db.Ping()).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Взаимодействие с таблицами БД", func() {
		It("При вставке значений в таблицу не должно быть ошибок", func() {
			stmt, err := db.Prepare("INSERT INTO table_executed_commands (id, service, binary_data) values(?,?,?)")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = stmt.Exec("d7873r734847", "MISP", []byte(`{[
			{
			  Command: "addtag",
			  String:  "Webhook: send=\"MISP\""
			}
			]}`))
			Expect(err).ShouldNot(HaveOccurred())
		})
		/*
			It("", func ()  {

			})
		*/
	})
})
