package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	// create key
	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestUpdate(t *testing.T) {
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

	// update key
	if err := ctx.Client.Update(ctx.Key1, ctx.Value2, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	// read update
	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value2, v)
	}
}

func TestDelete(t *testing.T) {
	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	assert := assert.New(t)

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, true)
	}

	if err := ctx.Client.Delete(ctx.Key1, TestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, false)
	}
}
