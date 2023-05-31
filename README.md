# kubectl-resource_report

kubectl-resource_report does things.

## Installation

```shell
# homebrew
brew install stenic/tap/kubectl-resource_report

# gofish
gofish rig add https://github.com/stenic/fish-food
gofish install github.com/stenic/fish-food/kubectl-resource_report

# scoop
scoop bucket add kubectl-resource_report https://github.com/stenic/scoop-bucket.git
scoop install kubectl-resource_report

# go
go install github.com/stenic/kubectl-resource_report@latest

# docker
docker pull ghcr.io/stenic/kubectl-resource_report:latest

# dockerfile
COPY --from=ghcr.io/stenic/kubectl-resource_report:latest /kubectl-resource_report /usr/local/bin/
```

> For even more options, check the [releases page](https://github.com/stenic/kubectl-resource_report/releases).

## Run

```shell
# Installed
kubectl-resource_report -h

# Docker
docker run -ti ghcr.io/stenic/kubectl-resource_report:latest -h

# Kubernetes
kubectl run kubectl-resource_report --image=ghcr.io/stenic/kubectl-resource_report:latest --restart=Never -ti --rm -- -h
```

## Documentation

```shell
kubectl-resource_report -h
```

## Badges

[![Release](https://img.shields.io/github/release/stenic/kubectl-resource_report.svg?style=for-the-badge)](https://github.com/stenic/kubectl-resource_report/releases/latest)
[![Software License](https://img.shields.io/github/license/stenic/kubectl-resource_report?style=for-the-badge)](./LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/stenic/kubectl-resource_report/Release?style=for-the-badge)](https://github.com/stenic/kubectl-resource_report/actions?workflow=build)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## License

[License](./LICENSE)
