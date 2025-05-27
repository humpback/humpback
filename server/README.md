# humpback-server

[![PkgGoDev](https://pkg.go.dev/badge/github.com/docker/docker)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/vue-3.x-brightgreen.svg)](https://v3.vuejs.org/)
[![Docker](https://img.shields.io/badge/docker-pull-blue?logo=docker)](https://hub.docker.com/r/humpbacks/humpback-server)
[![Release](https://img.shields.io/badge/release-v2.0.0-blue)](https://github.com/humpback/humpback-server/releases/tag/v2.0.0)

![Humpback logo](/assets/logo.png)

Lightweight container service management platform site.

## Language

- [English](README.md)
- [中文](README.zh.md)

## Feature

- Multiple Deployment Strategies: Flexible deployment strategies to meet your various business scenarios.
- Supporting Multiple Cluster：One-stop operation and management of multiple clusters.
- Centralized Access Control： Granular permission control (team and individual levels).
- Friendly Web UI: An intuitive web interface that hides the complexity of container operations.

## Getting Started

* [Humpback Guides](https://humpback.github.io/humpback)

## Installing

First, create the volume that Humpback Server will use to store its database:

```bash
docker volume create humpback_data
```

Then, install the Humpback Server container:

```bash
docker run -d \
  --name humpback-server \
  -p 8100:8100 \
  -p 8101:8101 \
  --restart=always \
  -v humpback_data:/workspace/data \
  -e LOCATION=prd \
  humpbacks/humpback-server
```

By default, Humpback Server will expose the UI over port `8100` and expose a API server over port `8101` for receiving
agent report.

Humpback Server has now been installed. you can log into your Humpback Server instance by opening a web browser and
going to:

```
http://localhost:8100
```

The initial super administrator account and password are **humpback**

## License

Humpback Server is licensed under the [Apache Licence 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).   