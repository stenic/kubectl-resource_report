# [[project]]

[[project]] does things.

## Installation

```shell
# homebrew
brew install stenic/tap/[[project]]

# gofish
gofish rig add https://github.com/stenic/fish-food
gofish install github.com/stenic/fish-food/[[project]]

# scoop
scoop bucket add [[project]] https://github.com/stenic/scoop-bucket.git
scoop install [[project]]

# go
go install github.com/stenic/[[project]]@latest

# docker 
docker pull ghcr.io/stenic/[[project]]:latest

# dockerfile
COPY --from=ghcr.io/stenic/[[project]]:latest /[[project]] /usr/local/bin/
```

> For even more options, check the [releases page](https://github.com/stenic/[[project]]/releases).


## Run

```shell
# Installed
[[project]] -h

# Docker
docker run -ti ghcr.io/stenic/[[project]]:latest -h

# Kubernetes
kubectl run [[project]] --image=ghcr.io/stenic/[[project]]:latest --restart=Never -ti --rm -- -h
```

## Documentation

```shell
[[project]] -h
```

## Badges

[![Release](https://img.shields.io/github/release/stenic/[[project]].svg?style=for-the-badge)](https://github.com/stenic/[[project]]/releases/latest)
[![Software License](https://img.shields.io/github/license/stenic/[[project]]?style=for-the-badge)](./LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/stenic/[[project]]/Release?style=for-the-badge)](https://github.com/stenic/[[project]]/actions?workflow=build)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## License

[License](./LICENSE)
