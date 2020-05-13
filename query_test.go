package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRead(t *testing.T) {
	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	assert := assert.New(t)

	// create key
	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	// read key
	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value1, v)
	}
}

//

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

//
