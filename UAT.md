### User Acceptance Testing

The following guide describes setting up the project and running an example code and tests in an AWS Ubuntu 18.04 machine. Once ssh'ed into the machine:

1. Ensure the system package index is up to date:

```
sudo apt -y update
```

2. Install required system tools

```
sudo apt install -y build-essential make
```

3. Ensure latest golang version is installed:

```
cd /tmp
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz

export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

4. Clone the project:

```
mkdir -p ~/go/src/github.com/bluzelle
cd  ~/go/src/github.com/bluzelle
git clone https://github.com/bluzelle/blzgo.git
cd blzgo
```

5. Setup the sample environment variables:

```
cp .env.sample .env
```

The example code and tests will read the bluzelle settings to use from that file i.e. `.env`.

6. Run the example code located at `examples/hello_world/main.go`:

```
make hello_world
```

This example code performs simple CRUD operations against the testnet.

7. The project also ships a complete suite of integration tests for all the methods. To run all the tests simply run:

```
make test
```

This will run all the tests in the `test` directory using the same environment settings defined in the `.env` file. A successful run should result in an output like this:

```
PASS
ok  	github.com/bluzelle/blzgo	231.143s
```
