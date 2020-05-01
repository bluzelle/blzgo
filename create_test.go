package bluzelle

import (
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	// create key
	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, 0, nil); err != nil {
		t.Fatalf("%s", err)
	}
}
