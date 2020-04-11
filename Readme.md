![](https://raw.githubusercontent.com/bluzelle/api/master/source/images/Bluzelle%20-%20Logo%20-%20Big%20-%20Colour.png)

### Getting started

Ensure you have a recent version of go installed.

Grab the package from github:

    $ go get -u github.com/vbstreetz/blzgo/src

Use:

```go
package main

import (
  "github.com/vbstreetz/blzgo"
)

// create client
options := &bluzelle.Options{
  Address:  "...",
  Mnemonic: "...",
  Endpoint: "http://testnet.public.bluzelle.com:1317",
}
client, err := bluzelle.NewClient(options)
if err != nil {
  log.Fatalf("%s", err)
}

// read account
if account, err := client.ReadAccount(); err != nil {
  log.Fatalf("%s", err)
} else {
  log.Printf("account info: %+v", account)
}
```

### Examples

You can test out the `examples/` included by:

1. Copy `.env.sample` to `.env` to configure the environment. Update the file `.env` accordingly, you can use the settings documented [here](https://docs.bluzelle.com/developers/bluzelle-db/getting-started-with-testnet).

2. Install dependencies:

    $ go get ./...

3. Run an example as defined in the `Makefile`, for example, to read the value of a previously set key(`foo`), run:

    $ make read key=foo

### Licence

MIT
