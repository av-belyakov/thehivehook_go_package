package testwebhooktemporarystorage_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	temporarystorage "github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver/temporarystorage"
)

var _ = Describe("Testwebhooktemporarystorage", Ordered, func() {
	var (
		whts *temporarystorage.WebHookTemporaryStorage

		err error

		test_uniq_case_id_1 string = "uniq_case_id:f78773r88r8w874et7rt7g77sw7w"
		test_uniq_case_id_2 string = "uniq_case_id:g7627gdff8fyr8298euihusd8y823"
		test_uniq_case_id_3 string = "uniq_case_id:fs662te73t73tr73t6rt37tr7376r3"

		test_uuid_1 string
		test_uuid_2 string
		test_uuid_3 string
	)

	_ = BeforeAll(func() {
		whts, err = temporarystorage.NewWebHookTemporaryStorage(10)
	})

	Context("Тест 1. Проверка работы webHookTemporaryStorage", func() {
		It("При инициализации модуля не должно быть ошибок", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("Новая информация должна быть успешно добавлена", func() {
			test_uuid_1 = whts.SetElementId(test_uniq_case_id_1)
			d, ok := whts.GetElementId(test_uuid_1)
			Expect(ok).Should(BeTrue())
			Expect(d).Should(Equal(test_uniq_case_id_1))

			test_uuid_2 = whts.SetElementId(test_uniq_case_id_2)
			d, ok = whts.GetElementId(test_uuid_2)
			Expect(ok).Should(BeTrue())
			Expect(d).Should(Equal(test_uniq_case_id_2))
		})
		It("Информация должна быть успешно удалена по её uuid", func() {
			whts.DeleteElement(test_uuid_1)
			_, ok := whts.GetElementId(test_uuid_1)
			Expect(ok).ShouldNot(BeTrue())
		})
		It("Информация должна быть успешно удалена по истечении её времени жизни", func() {
			time.Sleep(9 * time.Second)

			test_uuid_3 = whts.SetElementId(test_uniq_case_id_3)

			time.Sleep(6 * time.Second)
			//удаляется автоматически
			_, ok := whts.GetElementId(test_uuid_2)
			Expect(ok).ShouldNot(BeTrue())

			d, ok := whts.GetElementId(test_uuid_3)
			Expect(ok).Should(BeTrue())
			Expect(d).Should(Equal(test_uniq_case_id_3))
		})
	})

	/*
		Context("", func ()  {
			It("", func ()  {

			})
		})
	*/
})
