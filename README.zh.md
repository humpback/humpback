# Humpback 2.0

[![PkgGoDev](https://pkg.go.dev/badge/github.com/docker/docker)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/vue-3.x-brightgreen.svg)](https://v3.vuejs.org/)
[![Docker](https://img.shields.io/badge/docker-pull-blue?logo=docker)](https://hub.docker.com/r/humpbacks)
[![Release](https://img.shields.io/badge/release-v2.0.0-blue)](https://github.com/humpback/humpback/releases/tag/v2.0.0)

![Humpback logo](/docs/_media/logo-nobg.png)

轻量级容器化服务管理平台。

## 语言

- [English](README.md)
- [中文](README.zh.md)

## 文档

* [https://humpback.github.io/humpback](https://humpback.github.io/humpback)

## 架构

![Humpback Architecture](/docs/_media/humpback-architect-new.png)

## 特征

- 多种部署策略: 灵活多样的部署策略满足你的各种业务场景需求。
- 支持多集群：一站式运维多个集群。
- 集中式访问控制： 不同粒度的权限控制（团队和个人）。
- 简洁易用的界面: 简洁易用但功能强大的用户界面让你不需要面对容器的复杂性。

## 组件

* [Humpback Agent](https://github.com/humpback/humpback-agent)

## 快速开始

### 安装

首先，创建一个volume用于存储Humpback的数据库：

```bash
docker volume create humpback_data
```

接下，使用下面的命令创建Humpback容器：

```bash
docker run -d \
  --name humpback \
  --net=host \
  --restart=always \
  -v humpback_data:/workspace/data \
  -v humpback_certs:/workspace/certs \
  humpbacks/humpback
```

Humpback默认会监听两个端口，`8100`端口是web站点，`8101`是API服务器，主要接受agent汇报的数据。

命令运行成功后，你可以通过打开下面的站点检查Humpback是否启动成功。

```
http://localhost:8100
```

初始化的超级管理员账号密码均为 **humpback**

## 帮助

- 问题: https://github.com/humpback/humpback/issues
- Slack (chat): todo
- 微信： todo

## 许可证

Humpback 根据 [Apache Licence 2.0](http://www.apache.org/licenses/LICENSE-2.0.html) 获得许可。
