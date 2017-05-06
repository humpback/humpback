# 部署 Humpback Center 

> Humpback Center  

&ensp;&ensp;&ensp;Humpback Center 主要为 Humpback 平台提供集群容器调度服务，以集群中心角色实现各个 Group 的容器分配管理。   

&ensp;&ensp;&ensp;需要集群调度的 `Group` 要在设置中将 `Cluster Mode` 开关打开即可。

- 启动 Humpback Center

&ensp;&ensp;&ensp; `CENTER_LISTEN_PORT=:8589` Humpback Center 的 API 默认端口为：8589。   

&ensp;&ensp;&ensp; `DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181` 为先前配置的 `Zookeeper` 集群地址信息。

```bash 
$ docker pull humpback/humpback-center:1.0.0
$ docker run -d -ti --net=host --restart=always \
> --name=humpback-center \
> -e CENTER_LISTEN_PORT=:8589
> -e DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181 \
> -e DOCKER_CLUSTER_NAME=humpback/center \
> -e /opt/app/humpback-center/cache:/opt/humpback-center/cache \
> -e /opt/app/humpback-center/logs:/opt/humpback-center/logs \
> -v /var/run/:/var/run/:rw \
> humpbacks/humpback-center:1.0.0
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND                  CREATED        STATUS         PORTS         NAMES
a1640bf8c956    humpback/humpback-center:1.0.0  "/usr/bin/dumb-init -"   15 minutes ago  45 seconds ago              humpback-center
```

- [集群模式如何创建容器](cluster-create-container.md)

- [集群调度策略](cluster-container-schedule.md)