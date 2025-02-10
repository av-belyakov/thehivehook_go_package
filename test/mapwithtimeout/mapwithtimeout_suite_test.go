package mapwithtimeout_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMapwithtimeout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mapwithtimeout Suite")
}
