# kubectl-resource-report

kubectl-resource-report does things.

## Installation

```shell
# homebrew
brew install stenic/tap/kubectl-resource-report

# gofish
gofish rig add https://github.com/stenic/fish-food
gofish install github.com/stenic/fish-food/kubectl-resource-report

# scoop
scoop bucket add kubectl-resource-report https://github.com/stenic/scoop-bucket.git
scoop install kubectl-resource-report

# go
go install github.com/stenic/kubectl-resource-report@latest

# docker
docker pull ghcr.io/stenic/kubectl-resource-report:latest

# dockerfile
COPY --from=ghcr.io/stenic/kubectl-resource-report:latest /kubectl-resource-report /usr/local/bin/
```

> For even more options, check the [releases page](https://github.com/stenic/kubectl-resource-report/releases).

## Run

```shell
# Installed
kubectl-resource-report -h

# Docker
docker run -ti ghcr.io/stenic/kubectl-resource-report:latest -h

# Kubernetes
kubectl run kubectl-resource-report --image=ghcr.io/stenic/kubectl-resource-report:latest --restart=Never -ti --rm -- -h
```

## Documentation

```shell
kubectl-resource-report -h
```

## Badges

[![Release](https://img.shields.io/github/release/stenic/kubectl-resource-report.svg?style=for-the-badge)](https://github.com/stenic/kubectl-resource-report/releases/latest)
[![Software License](https://img.shields.io/github/license/stenic/kubectl-resource-report?style=for-the-badge)](./LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/stenic/kubectl-resource-report/Release?style=for-the-badge)](https://github.com/stenic/kubectl-resource-report/actions?workflow=build)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## License

[License](./LICENSE)
