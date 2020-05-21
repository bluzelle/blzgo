package bluzelle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if account, err := ctx.Client.Account(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(account.AccountNumber > 0)
		assert.True(account.Sequence > 0)
	}
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if v, err := ctx.Client.Version(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v != "")
	}
}

//

func TestRead(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value1, v)
	}
}

func TestHas(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if v, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(!v)
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v)
	}
}

func TestCount(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	count, err := ctx.Client.Count()
	if err != nil {
		t.Fatalf("%s", err)
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if newAccount, err := ctx.Client.Count(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(newAccount, count+1)
	}
}

func TestKeys(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if keys, err := ctx.Client.Keys(); err != nil {
		t.Fatalf("%s", err)
	} else {
		for _, k := range keys {
			assert.True(k != ctx.Key1)
		}
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if keys, err := ctx.Client.Keys(); err != nil {
		t.Fatalf("%s", err)
	} else {
		found := false
		for _, k := range keys {
			if k == ctx.Key1 {
				found = true
			}
		}
		if !found {
			t.Fatalf("key(%s) was not found in (%+v)", ctx.Key1, keys)
		}
	}
}

func TestKeyValues(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if kvs, err := ctx.Client.KeyValues(); err != nil {
		t.Fatalf("%s", err)
	} else {
		for _, kv := range kvs {
			assert.True(kv.Key != ctx.Key1)
		}
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if kvs, err := ctx.Client.KeyValues(); err != nil {
		t.Fatalf("%s", err)
	} else {
		found := false
		for _, kv := range kvs {
			if kv.Key == ctx.Key1 {
				found = true
			}
		}
		if !found {
			t.Fatalf("key(%s) was not found in (%+v)", ctx.Key1, kvs)
		}
	}
}

func TestGetLease(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.GetLease(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v > 0)
	}
}

func TestGetNShortestLeases(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, TestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.GetNShortestLeases(10); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(len(v) > 0)
	}
}
