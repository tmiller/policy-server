# policy-server

This is the readme for the golang policy server. Please see development setup
for things you will need to develop or deploy this project

### Development Setup

#### Setup Build tools (required for ruby version too)

```bash
brew install gnu-tar
gem install fpm
```

#### Setup golang

Download and install go from [here](https://golang.org/dl/).
```bash
# put this in .zshrc or .bashrc
export PATH=$PATH:/usr/local/go/bin
```

#### Pull down repository (go assumes you put code in `~/go/src`)

```bash
mkdir -p ~/go/src/github.com/tmiller/policy-server
cd ~/go/src/github.com/tmiller/policy-server
```

### Building

```bash
# Build Release
git merge <feature-branch>
git tag -a <version-number> # make sure to look at other tags for tag format
VERSION=<version-number> make all
```

### Running the Policy Server

#### Create self-signed certifiate for TLS

To be able to run the policy server a self signed TLS certificate is required
to generate this certificate you can run the following. The only thing you need
to fill out is "Common Name" all other questions can be ignored by pressing
enter at each prompt

```bash
# Use 'localhost' for the 'Common name'
openssl req -x509 -sha256 -nodes -newkey rsa:2048 -days 365 -keyout localhost.key -out localhost.crt
```

#### Running the policy server

Running the policy-server using `go run`:

```bash
go run main.go -b :8080 -c localhost.crt -k localhost.key -p resources/crossdomain.xml
```

Running the policy-server using the built binary:

```bash
build/darwin-amd64/policy-server -b :8080 -c localhost.crt -k localhost.key -p resources/crossdomain.xml
```

#### Connecting to the policy server

You can use `curl` though it is technically not an HTTP server. You can also use
`socat` (`brew install socat`). Curl is probably easiest to use

```bash
# Using curl
curl --cacert localhost.crt https://localhost:8080

# Using socat
socat ssl:localhost:8080,cafile=localhost.crt stdio
# Press enter
```

### Performance Tuning

There are two flags that can be used to tune performance. These are the number
of workers handling connections and connection queue size.

#### Workers

*Valid Values*: You can safely pass one and greater to this flag, it defaults
to one.

Configuring the number of workers directly affects CPU load. The more workers
running the more a CPU will be taxed.

#### Connection Queue

*Valid Values*: You can safely pass zero or greater to this flag, it defaults
to zero.

*Warning* setting the queue too high can cause the go process to run out of
open files. On a default linux configuration which has an allowed 1024 open
files per process, a safe number is around 900 (allowing head room for workers
to hold open connections).

Configuring the connection queue can help provide consistent performance in
regards to timeouts. If this value is 0 there is no queue and the workers block
until there is a request to read. When there is a lot of load this can cause
timeouts to happen to a number of requests. When this value is one or greater,
then a queue is setup to hold open connections until they can be responded to.
This causes latency but can provided consistent performance can cause less
timeouts.
