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
