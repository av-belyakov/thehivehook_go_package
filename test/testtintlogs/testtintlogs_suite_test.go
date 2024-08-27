package testtintlogs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTesttintlogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testtintlogs Suite")
}
