package bluzelle

import (
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := &Test{}
	if err := ctx.SetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TearDown()

	// create key
	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, 0); err != nil {
		t.Fatalf("%s", err)
	}
}
