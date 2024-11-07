package nats_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestnats(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testnats Suite")
}
