package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequiresMnemonic(t *testing.T) {
	assert := assert.New(t)

	options := &Options{}
	if _, err := NewClient(options); err == nil {
		t.Fatalf("mnemonic requirement was not validated")
	} else {
		assert.Equal(err.Error(), "mnemonic is required")
	}
}

func TestRequiresUUID(t *testing.T) {
	assert := assert.New(t)

	options := &Options{
		Mnemonic: "...",
	}
	if _, err := NewClient(options); err == nil {
		t.Fatalf("UUID requirement was not validated")
	} else {
		assert.Equal(err.Error(), "uuid is required")
	}
}

func TestAddressDerivation(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	// compute addr
	if err := ctx.Client.setAddress(); err != nil {
		t.Fatalf("%s", err)
	}

	// validate addr
	assert.Equal(ctx.Client.Address, TestAddress())
}

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
