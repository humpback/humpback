# 安装Humpback

## 先决条件

Humpback 由两个核心组件构成：Humpback Server 和 Humpback Agent。这两个组件都以轻量级 Docker 容器的形式运行在 Docker 引擎上。

开始使用前，您需要安装最新版本的 Docker 并确保其正常运行。我们建议您遵循 Docker 的[官方安装指南](https://docs.docker.com/engine/install/)——特别提醒，在 Ubuntu 系统上不建议通过 snap 安装 Docker，否则可能会遇到兼容性问题。

## 部署Hummpback Server

首先，创建一个volume用于存储Humpback Server的数据库：

```bash
docker volume create humpback_data
```

接下，使用下面的命令创建Humpback Server容器：

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

Humpback Server默认会监听两个端口，`8100`端口是web站点，`8101`是API服务器，主要接受agent汇报的数据。

命令运行成功后，你可以通过打开下面的站点检查Humpback Server是否启动成功。

```
http://localhost:8100
```

你可以使用系统初始化的用户登录，用户名和密码都是 `humpback`。

## 部署Hummpback Agent

Humpback Agent默认会监听`8018`端口用于接收Humpback Server的调用。

使用下面的命令启动Humpback Agent

```bash

docker run -d \
  --name=humpback-agent \
  --net=host \
  --restart=always \
  --privileged \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/lib/docker:/var/lib/docker \
  -e HUMPBACK_AGENT_API_BIND=0.0.0.0:8018 \
  -e HUMPBACK_SERVER_HOST={server-address}:8101 \
  -e HUMPBACK_VOLUMES_ROOT_DIRECTORY=/var/lib/docker \
  humpbacks/humpback-agent

```
请注意：将`{server-address}`替换为部署Humpback Server的真实IP地址。

## 在Humpback中添加节点

打开Humpback网站，在Administration菜单下，打开Nodes页面，你可以通过IP地址添加刚刚安装Agent的节点。