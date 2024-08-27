package testcontexts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestcontexts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testcontexts Suite")
}
