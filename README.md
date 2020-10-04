![](https://raw.githubusercontent.com/bluzelle/api/master/source/images/Bluzelle%20-%20Logo%20-%20Big%20-%20Colour.png)

### Getting started

Ensure you have a recent version of [Go](https://golang.org) installed.

Grab the package from github:

```
go get github.com/bluzelle/blzgo
```

Use:

```go
package main

import (
  "github.com/bluzelle/blzgo"
  "log"
)

func main() {
  // create client
  options := &bluzelle.Options{
    Mnemonic: "...",
    Endpoint: "http://dev.testnet.public.bluzelle.com:1317",
    UUID: "...",
  }
  client, err := bluzelle.NewClient(options)
  if err != nil {
    log.Fatalf("%s", err)
  }

  key := "foo"
  value := "bar"

  gasInfo := bluzelle.GasInfo{
		MaxFee:   4000001,
		MaxGas:   200000,
		GasPrice: 10,
  }
  leaseInfo := bluzelle.LeaseInfo{
    Days: 1,
  }

  // create key
  if err := client.Create(key, value, gasInfo, leaseInfo); err != nil {
    log.Fatalf("%s", err)
  } else {
    log.Printf("create key success: true\n")
  }

  // read key
  if v, err := client.Read(key); err != nil {
    log.Fatalf("%s", err)
  } else {
    log.Printf("read key success: %t\n", v == value)
  }

  // update key
  if err := client.Create(key, value, gasInfo, nil); err != nil {
    log.Fatalf("%s", err)
  } else {
    log.Printf("create key success: true\n")
  }

  // delete key
  if err := client.Delete(key, gasInfo); err != nil {
    log.Fatalf("%s", err)
  } else {
    log.Printf("delete key success: true\n")
  }
}
```

### Examples

You can test out the `examples/` included by following these steps:

1. Copy `.env.sample` to `.env` and configure if needed.

```
cp .env.sample .env
```

2. Install dependencies:

```
go get ./...
```

3. Run an example as defined in the `Makefile`, for example, to read the value of an existing key, `foo`, run:

```
make read key=foo
```

This will run the `examples/crud/read.go`.

### Integration Tests

```
make test
```

### User Acceptance Testing

Please checkout the [UAT.md](https://github.com/bluzelle/blzgo/blob/master/UAT.md) document for more details.

### Licence

MIT
