# policy-server

This is the readme for the golang policy server. Please see development setup
for things you will need to develop or deploy this project

### Development Setup

```
brew install gnu-tar
gem install fpm
```

### Building

```
# Build Release
git merge <feature-branch>
git tag -a <version-number> # make sure to look at other tags for tag format
make all
```
