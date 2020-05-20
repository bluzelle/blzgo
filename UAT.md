### User Acceptance Testing

The following guide describe setting up the project and running an example code and tests in an Ubuntu 18.04 machine. Once ssh'd into the machine:

1. Ensure latest go version is installed:

```
sudo apt install -y golang
```

2. Setup the sample environment variables:

```
cp .env.sample .env
```

The example code and tests will read the bluzelle settings to use from that file i.e. `.env`.

7. Run the example code located at `examples/hello_world/main.go`:

```
make hello_world
```

This example code performs simple CRUD operations against the testnet.

8. The project also ships a complete suite of integration tests for all the methods. To run all the tests simply run:

```
make test
```

This will run all the tests in the `test` directory using the same environment settings defined in the `.env` file.
Note that sometimes one or 2 tests fail due to some existing issues with the testnet. A successful run should result in an output like this:

```
PASS
ok  	github.com/vbstreetz/blzgo	34.302s
```
