package testcontexts_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testcontexts", func() {
	Context("Тест 1. Проверка обработки нескольких context", func() {
		It("Должен быть успешно обработан context", func() {
			ctx, cancel := context.WithCancel(context.Background())

			//ctxa, _ := context.WithCancel(ctx)
			go func() {
				fmt.Println("111 any groutina")

				<-ctx.Done()

				fmt.Println("111 Close second context")
			}()

			go func() {
				fmt.Println("222 any groutina")

				<-ctx.Done()

				fmt.Println("222 Close second context")
			}()

			time.Sleep(1 * time.Second)

			cancel()
			fmt.Println("Close first context")

			time.Sleep(2 * time.Second)

			Expect(true).Should(BeTrue())
		})
	})

	Context("Тест 2. Тестируем NotifyContext на нескольких гроутинах", func() {
		It("Должны быть успешно остановлены три вспомогательные гроутины и одна основная", func() {
			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT)

			go func() {
				sigChan := make(chan os.Signal, 1)
				osCall := <-sigChan
				log.Printf("system call:%+v", osCall)

				cancel()
			}()

			var wg sync.WaitGroup
			wg.Add(3)

			go func(myctx context.Context) {
				fmt.Println("function №1 start...")

				<-myctx.Done()
				fmt.Println("function №1 end")
				wg.Done()
			}(ctx)

			go func(myctx context.Context) {
				fmt.Println("function №2 start...")

				<-myctx.Done()
				fmt.Println("function №2 end")
				wg.Done()
			}(ctx)

			go func(myctx context.Context) {
				fmt.Println("function №3 start...")

				<-myctx.Done()
				fmt.Println("function №3 end")
				wg.Done()
			}(ctx)

			go func(cancel context.CancelFunc) {
				time.Sleep(3 * time.Second)
				cancel()
			}(cancel)

			wg.Wait()

			fmt.Println("Test STOP ")

			Expect(true).Should(BeTrue())
		})
	})
})
