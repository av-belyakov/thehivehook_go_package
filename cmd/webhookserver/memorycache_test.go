package webhookserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/v2/cmd/webhookserver"
)

func TestMemCache(t *testing.T) {
	mcOne := webhookserver.NewMemoryCache()
	mcOne.Alloc = 624624
	assert.Equal(t, mcOne.Alloc, uint64(624624))

	mcTwo := webhookserver.NewMemoryCache()

	assert.Equal(t, mcTwo.Alloc, uint64(624624))
	mcTwo.Alloc = 6742845
	assert.Equal(t, mcTwo.Alloc, uint64(6742845))
}
