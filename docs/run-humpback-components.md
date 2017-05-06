# 如何启动Humpback

> 准备工作   

首先准备好服务器，以下三台机器为搭建集群示例，在开始前请为每台机器安装好 `Docker` 实例。   

```
SERVER01：192.168.2.80 
SERVER02：192.168.2.81
SERVER03：192.168.2.82 
```

> 部署介绍

&ensp;&ensp;&ensp;可以看到 `192.168.2.80` 作为 Humpback 平台中心，启动的服务最多：   

&ensp;&ensp;&ensp;Docker Registry，镜像私有仓库服务。    

&ensp;&ensp;&ensp;Humpback-Web，Humpback 管理站点。      

&ensp;&ensp;&ensp;Humpback-Center，为集群模式提供容器调度服务。   

&ensp;&ensp;&ensp;同时以下三台服务器都启动 `Humpback-Agent` 提供本地镜像管理功能，等待被 Humpback 平台调用。   

```
| Server           |  Zookeeper |  Docker Registry  |  Humpback-Web  |  Humpback-Agent  |  Humpback-Center  |
|------------------|:----------:|:-----------------:|:--------------:|:----------------:|:-----------------:|
|   192.168.2.80   |     √      |         √         |        √       |        √         |         √         |
|   192.168.2.81   |     √      |         X         |        X       |        √         |         X         |
|   192.168.2.82   |     √      |         X         |        X       |        √         |         X         |
```

&ensp;&ensp;&ensp;每台服务器都需要启动 `Zookeeper` 节点组成集群，但搭建 `Zookeeper` 集群不是必选项。   
  
&ensp;&ensp;&ensp;Humpback 平台管理模式分为两种：   

&ensp;&ensp;&ensp;`Single Mode`，在该模式下 `Humpback-Web` 会直接访问 `Humpback-Agent` API 进行交互实现简单的集群容器管理，这种模式部署更单一，不必启动 `Humpback-Center` 和搭建 `Zookeeper` 集群。   

&ensp;&ensp;&ensp;`Cluster Mode`，Humpack 后端需要启动 `Humpback-Center` 和 `Zookeeper` 集群来支撑，`Humpback-Center` 作为集群调度中心提供容器动态调度分配功能，而 `Zookeeper` 集群提供节点发现功能，实现集群节点弹性伸缩。

> Zookeeper 集群

&ensp;&ensp;&ensp;[如何 Zookeeper 集群搭建](run-zookeeper.md) 

> Docker Registry 私有仓库

&ensp;&ensp;&ensp;[如何部署 Docker Registry](run-registry.md)

> Humpback Web 站点

&ensp;&ensp;&ensp;[如何部署 Humpback 站点](run-humpback-web.md)

> Humpback Agent

&ensp;&ensp;&ensp;[如何部署 Humpback Agent](run-humpback-agent.md)

> Humpback Center

&ensp;&ensp;&ensp;[如何部署 Humpback Center](run-humpback-center.md)









 