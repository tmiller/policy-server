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
Put this in `~/.bashrc` or `~/.zshrc`
```bash
export PATH=$PATH:/usr/local/go/bin
```

#### Pull down repository (go assumes you put code in `~/go/src`)

```bash
mkdir -p ~/go/src/github.com/tmiller/policy-server
cd ~/go/src/github.com/tmiller/policy-server
```

### Building

```
# Build Release
git merge <feature-branch>
git tag -a <version-number> # make sure to look at other tags for tag format
make all
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
