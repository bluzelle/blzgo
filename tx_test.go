package bluzelle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWithGas(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}
	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(v, ctx.Value1)
	}
}

func TestCreateWithLeaseInfo(t *testing.T) {
	// assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), &LeaseInfo{Seconds: 60}); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateValidatesGasInfo(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	err := ctx.Client.Create(ctx.Key1, ctx.Value1, &GasInfo{MaxFee: 1}, nil)
	assert.True(err != nil) // todo check details
}

func TestCreateKeyWithSymbols(t *testing.T) {
	//assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	key := ctx.Key1 + " !\"#$%&'()*+,-.0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	val := ctx.Value1

	if err := ctx.Client.Create(key, val, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	// if v, err := ctx.Client.Read(key); err != nil {
	// 	t.Fatalf("%s", err)
	// } else {
	// 	assert.Equal(v, val)
	// }

	if keys, err := ctx.Client.Keys(); err != nil {
		t.Fatalf("%s", err)
	} else {
		found := false
		// t.Logf("%+v", keys)
		for _, k := range keys {
			if k == key {
				found = true
			}
		}
		if !found {
			t.Fatalf("key(%s) was not found in (%+v)", ctx.Key1, keys)
		}
	}
}

func TestCreatesFailsIfKeyContainsSlash(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	err := ctx.Client.Create("123/", ctx.Value1, GetTestGasInfo(), nil)
	assert.True(err != nil)
	assert.True(strings.Contains(err.Error(), "Key cannot contain a slash"))
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if err := ctx.Client.Update(ctx.Key1, ctx.Value2, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.Read(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value2, v)
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, true)
	}

	if err := ctx.Client.Delete(ctx.Key1, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(b, false)
	}
}

func TestRename(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key2); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(!b)
	}

	if err := ctx.Client.Rename(ctx.Key1, ctx.Key2, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	}

	if b, err := ctx.Client.Has(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(!b)
	}

	if b, err := ctx.Client.Has(ctx.Key2); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(b)
	}
}

func TestDeleteAll(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if count, err := ctx.Client.Count(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(count >= 1)
	}

	if err := ctx.Client.DeleteAll(GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	}

	if count, err := ctx.Client.Count(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(0, count)
	}
}

func TestMultiUpdate(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.KeyValues(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(len(v), 1)
		assert.Equal(v[0].Key, ctx.Key1)
		assert.Equal(v[0].Value, ctx.Value1)
	}

	kvs := []*KeyValue{
		&KeyValue{ctx.Key1, "hey"},
	}

	if err := ctx.Client.MultiUpdate(kvs, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.KeyValues(); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(len(v), 1)
		assert.Equal(v[0].Key, ctx.Key1)
		assert.Equal(v[0].Value, "hey")
	}
}

func TestRenewKeyLease(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if lease, err := ctx.Client.GetLease(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(lease > 20)
	}

	if err := ctx.Client.RenewLease(ctx.Key1, GetTestGasInfo(), &LeaseInfo{Seconds: 10}); err != nil {
		t.Fatalf("%s", err)
	}

	if lease, err := ctx.Client.GetLease(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(lease < 20)
	}
}

func TestRenewAllKeyLeases(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if lease, err := ctx.Client.GetLease(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(lease > 20)
	}

	if err := ctx.Client.RenewAllLeases(GetTestGasInfo(), &LeaseInfo{Seconds: 10}); err != nil {
		t.Fatalf("%s", err)
	}

	if lease, err := ctx.Client.GetLease(ctx.Key1); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(lease < 20)
	}
}

//

func TestTxRead(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.TxRead(ctx.Key1, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(ctx.Value1, v)
	}
}

func TestTxHas(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if v, err := ctx.Client.TxHas(ctx.Key1, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(!v)
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.TxHas(ctx.Key1, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v)
	}
}

func TestTxCount(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	count, err := ctx.Client.TxCount(GetTestGasInfo())
	if err != nil {
		t.Fatalf("%s", err)
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if newAccount, err := ctx.Client.TxCount(GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.Equal(newAccount, count+1)
	}
}

func TestTxKeys(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if keys, err := ctx.Client.TxKeys(GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		for _, k := range keys {
			assert.True(k != ctx.Key1)
		}
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if keys, err := ctx.Client.TxKeys(GetTestGasInfo()); err != nil {
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

func TestTxKeyValues(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if kvs, err := ctx.Client.TxKeyValues(GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		for _, kv := range kvs {
			assert.True(kv.Key != ctx.Key1)
		}
	}

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if kvs, err := ctx.Client.TxKeyValues(GetTestGasInfo()); err != nil {
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

func TestTxGetLease(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.TxGetLease(ctx.Key1, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(v > 0)
	}
}

func TestTxGetNShortestLeases(t *testing.T) {
	assert := assert.New(t)

	ctx := &Test{}
	if err := ctx.TestSetUp(); err != nil {
		t.Fatalf("%s", err)
	}
	defer ctx.TestTearDown()

	if err := ctx.Client.Create(ctx.Key1, ctx.Value1, GetTestGasInfo(), nil); err != nil {
		t.Fatalf("%s", err)
	}

	if v, err := ctx.Client.TxGetNShortestLeases(10, GetTestGasInfo()); err != nil {
		t.Fatalf("%s", err)
	} else {
		assert.True(len(v) > 0)
	}
}
