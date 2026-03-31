package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// New конструктор обработчиа логов (это просто мост соединяющий несколько сервисов логирования)
func New(writer interfaces.WriterLoggingData, chSysMonit chan<- interfaces.Messager) *LoggingChan {
	return &LoggingChan{
		dataWriter:           writer,
		chanSystemMonitoring: chSysMonit,
		chanLogging:          make(chan interfaces.Messager),
	}
}

// Start обработчик и распределитель логов
func (lc *LoggingChan) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			//основной канал для приёма логов
			case msg := <-lc.chanLogging:
				//**********************************************************************
				//здесь так же может быть вывод в консоль, с индикацией цветов соответствующих
				//определенному типу сообщений но для этого надо включить вывод на stdout
				//в конфигурационном файле
				lc.dataWriter.Write(msg.GetType(), msg.GetMessage())

				if msg.GetType() == "error" || msg.GetType() == "warning" {
					msg := NewMessageLogging()
					msg.SetType("error")
					msg.SetMessage(fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()))

					lc.chanSystemMonitoring <- msg
				}

				if msg.GetType() == "info" {
					msg := NewMessageLogging()
					msg.SetType("info")
					msg.SetMessage(msg.GetMessage())

					lc.chanSystemMonitoring <- msg
				}
			}
		}
	}()
}
