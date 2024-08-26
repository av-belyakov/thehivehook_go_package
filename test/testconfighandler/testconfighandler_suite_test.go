package testconfighandler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestconfighandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testconfighandler Suite")
}
