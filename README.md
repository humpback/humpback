# Humpback 2.0

[![PkgGoDev](https://pkg.go.dev/badge/github.com/docker/docker)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/vue-3.x-brightgreen.svg)](https://v3.vuejs.org/)
[![Docker](https://img.shields.io/badge/docker-pull-blue?logo=docker)](https://hub.docker.com/r/humpbacks)
[![Release](https://img.shields.io/badge/release-v2.0.0-blue)](https://github.com/humpback/humpback/releases/tag/v2.0.0)

![Humpback logo](/docs/_media/logo-nobg.png)

Lightweight platform for managing containerized services.

## Language

- [English](README.md)
- [中文](README.zh.md)

## Documents

* [https://humpback.github.io/humpback](https://humpback.github.io/humpback) 

## Architecture

![Humpback Architecture](/docs/_media/humpback-architect-new.png)

## Feature

- Multiple Deployment Strategies: Flexible deployment strategies to meet your various business scenarios.
- Supporting Multiple Cluster：One-stop operation and management of multiple clusters.
- Centralized Access Control： Granular permission control (team and individual levels).
- Friendly Web UI: An intuitive web interface that hides the complexity of container operations.

## Components

* [Humpback Agent](https://github.com/humpback/humpback-agent)

## Getting Started

### Installing

First, create the volume that Humpback will use to store its database:

```bash
docker volume create humpback_data
```

Then, install the Humpback container:

```bash
docker run -d \
  --name humpback \
  --net=host \
  --restart=always \
  -v humpback_data:/workspace/data \
  -v humpback_certs:/workspace/certs \
  humpbacks/humpback
```

By default, Humpback will expose the UI over port `8100` and expose a API server over port `8101` for receiving
agent report.

Humpback has now been installed. you can log into your Humpback instance by opening a web browser and
going to:

```
http://localhost:8100
```

The initial super administrator account and password are **humpback**


## Getting Help

- Issues: https://github.com/humpback/humpback/issues
- Slack (chat): todo
- Wechat： todo

## License

Humpback is licensed under the [Apache Licence 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).   
