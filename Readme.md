![](https://raw.githubusercontent.com/bluzelle/api/master/source/images/Bluzelle%20-%20Logo%20-%20Big%20-%20Colour.png)

### Getting started

Ensure you have a recent version of [Go](https://golang.org) installed.

Grab the package from github:

    $ go get github.com/vbstreetz/blzgo

Use:

```go
package main

import (
  "github.com/vbstreetz/blzgo"
  "log"
)

func main() {
  // create client
  options := &bluzelle.Options{
    Address:  "...",
    Mnemonic: "...",
    Endpoint: "http://testnet.public.bluzelle.com:1317",
    GasInfo: &bluzelle.GasInfo{
      MaxFee: 4000001,
    },
  }
  client, err := bluzelle.NewClient(options)
  if err != nil {
    log.Fatalf("%s", err)
  }

  key := "foo"
  value := "bar"

  // create key
  if err := client.Create(key, value); err != nil {
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

  // delete key
  if err := client.Delete(key); err != nil {
    log.Fatalf("%s", err)
  } else {
    log.Printf("delete key success: true\n")
  }
}
```

### Examples

You can test out the `examples/` included by:

1. Coping `.env.sample` file to `.env` to configure the environment and then updating the resulting file, `.env`, accordingly. You can also use this test [file](https://gist.github.com/vbstreetz/f05a982530311d155836e27d41c1f73a)

2. Install dependencies:
```
    $ go get ./...
```
3. Run an example as defined in the `Makefile`, for example, to read the value of an existing key, `foo`, run:
```
    $ make read key=foo
```
### Integration Tests
```
    $ make test
```
### Licence

MIT
