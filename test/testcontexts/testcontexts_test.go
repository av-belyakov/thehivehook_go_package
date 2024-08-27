package testcontexts_test

import (
	"context"
	"fmt"
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
})
