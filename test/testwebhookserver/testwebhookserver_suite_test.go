package testwebhookserver_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestwebhookserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testwebhookserver Suite")
}
