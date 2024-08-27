package testtintlogs_test

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testtintlogs", func() {
	Context("Тест 1. Проверка пакета для вывода логов", func() {
		It("Должны быть успешно выполнены выводы логов", func() {
			loggerColor := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.DateTime,
			}))

			loggerColor.Error("Error test")
			loggerColor.Warn("Warning test")
			loggerColor.Info("Info test")

			Expect(true).Should(BeTrue())
		})
	})
})
