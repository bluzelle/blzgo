package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNShortestLeasesHumanized(t *testing.T) {
	assetTestGetNShortestLeasesHumanized(t, "1", 5)
	assetTestGetNShortestLeasesHumanized(t, "2", 10)
}

func assetTestGetNShortestLeasesHumanized(t *testing.T, lease string, seconds int64) {
	assert := assert.New(t)

	result := &GetNShortestLeasesResponseResult{
		KeyLeases: []*KeyLease{
			&KeyLease{Lease: lease},
		},
	}
	if humanized, err := result.GetHumanizedKeyLeases(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(humanized[0].Lease, seconds)
	}
}
