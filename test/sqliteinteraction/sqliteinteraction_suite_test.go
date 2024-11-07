package sqliteinteraction_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestsqliteinteraction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testsqliteinteraction Suite")
}
