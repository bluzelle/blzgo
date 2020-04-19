package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	ctx := &Test{}
	if err := ctx.SetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TearDown()

	assert := assert.New(t)

	// create key
	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, 0); err != nil {
		t.Fatalf("%s", err)
	}

	// update key
	if err := ctx.Client.Update(ctx.Key1, ctx.Value2, 0); err != nil {
		t.Fatalf("%s", err)
	}

	// read update
	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value2, v)
	}
}
