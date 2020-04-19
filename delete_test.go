package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	ctx := &Test{}
	if err := ctx.SetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TearDown()

	assert := assert.New(t)

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, 0); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, true)
	}

	if err := ctx.Client.Delete(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, false)
	}
}
