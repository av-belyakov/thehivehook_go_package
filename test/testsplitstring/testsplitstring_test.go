package testsplitstring_test

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

var _ = Describe("Testsplitstring", func() {
	strExample := "caseupdate:gcm,rcmmsk,rcmnvs;alertupdate:gcm"
	subscribers := []confighandler.SubscriberNATS{}

	hundlerSubscribersString := func(str string) (confighandler.SubscriberNATS, error) {
		errMsg := "an incorrect string containing the 'subscribers' of the NATS settings was received"
		subscriber := confighandler.SubscriberNATS{}

		if !strings.Contains(str, ":") {
			return subscriber, errors.New(errMsg)
		}

		tmp := strings.Split(str, ":")
		if len(tmp) == 0 {
			return subscriber, errors.New(errMsg)
		}

		responders := []string{}
		if strings.Contains(tmp[1], ",") {
			responders = append(responders, strings.Split(tmp[1], ",")...)
		} else {
			responders = append(responders, tmp[1])
		}

		subscriber.Event = tmp[0]
		subscriber.Responders = responders

		return subscriber, nil
	}

	Context("Тест 1. Деление строки", func() {
		It("Тестовая строка должна быть успешно разделена", func() {
			if strings.Contains(strExample, ";") {
				tmp := strings.Split(strExample, ";")
				if len(tmp) > 0 {
					for _, v := range tmp {
						subscriber, err := hundlerSubscribersString(v)
						Expect(err).ShouldNot(HaveOccurred())

						subscribers = append(subscribers, subscriber)
					}
				}
			} else {
				subscriber, err := hundlerSubscribersString(strExample)
				Expect(err).ShouldNot(HaveOccurred())

				subscribers = append(subscribers, subscriber)
			}

			fmt.Println("-------------------")
			fmt.Println(subscribers)
			fmt.Println("-------------------")

			Expect(len(subscribers)).Should(Equal(2))

			var num int
			for _, v := range subscribers {
				if v.Event == "caseupdate" {
					num = len(v.Responders)
				}
			}

			Expect(num).Should(Equal(3))
		})
	})
})
