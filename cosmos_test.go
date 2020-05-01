package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLeaseInfoToBlocks(t *testing.T) {
	assert := assert.New(t)

	assert.Equal((&LeaseInfo{}).ToBlocks(), int64(0))
	assert.Equal((&LeaseInfo{Seconds: 5}).ToBlocks(), int64(1))
	assert.Equal((&LeaseInfo{Seconds: 5, Minutes: 1}).ToBlocks(), int64(13))
	assert.Equal((&LeaseInfo{Seconds: 5, Minutes: 1, Hours: 1}).ToBlocks(), int64(733))
	assert.Equal((&LeaseInfo{Seconds: 5, Minutes: 1, Hours: 1, Days: 1}).ToBlocks(), int64(18013))

}
