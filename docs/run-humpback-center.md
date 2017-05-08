# 部署 Humpback Center 

> Humpback Center  

&ensp;&ensp;&ensp;Humpback Center 主要为 Humpback 平台提供集群容器调度服务，以集群中心角色实现各个 Group 的容器分配管理。   

&ensp;&ensp;&ensp;需要集群调度的 `Group` 要在设置中将 `Cluster Mode` 开关打开即可。

- 启动 Humpback Center

```bash 
$ docker pull humpbacks/humpback-center:1.0.0
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-center \
 -e HUMPBACK_SITEAPI=http://192.168.2.80:8012/api \
 -e CENTER_LISTEN_PORT=:8589 \
 -e CENTER_API_ENABLECORS=true \
 -e DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181 \
 -e DOCKER_CLUSTER_NAME=humpback/center \
 -v /opt/app/humpback-center/cache:/opt/humpback-center/cache \
 -v /opt/app/humpback-center/logs:/opt/humpback-center/logs \
 humpbacks/humpback-center:1.0.0
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND                  CREATED         STATUS         PORTS         NAMES
a1640bf8c956    humpbacks/humpback-center:1.0.0  "./humpback-center"     15 minutes ago  45 seconds ago              humpback-center
```

- 环境变量与参数

&ensp;&ensp;&ensp;`HUMPBACK_SITEAPI=http://192.168.2.80:8012/api` Humpback-Web 站点地址，注意要带上 `/api`。   

&ensp;&ensp;&ensp;`CENTER_LISTEN_PORT=:8589` Humpback Center 的 API 默认端口为：8589。   

&ensp;&ensp;&ensp;`CENTER_API_ENABLECORS=true` Humpback Center API 是否支持跨域访问。

&ensp;&ensp;&ensp;`DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181` 为先前配置的 `Zookeeper` 集群地址信息。

&ensp;&ensp;&ensp;`DOCKER_CLUSTER_NAME=humpback/center` 集群名称，要与 `Humpback Agnet` 配置一致。   

&ensp;&ensp;&ensp;`-v /opt/app/humpback-center/cache` 集群容器信息持久化目录，建议不要手动更改与删除。   

&ensp;&ensp;&ensp;`-v /opt/app/humpback-center/logs` 系统日志目录。

- [集群模式如何创建容器](cluster-create-container.md)

- [集群调度策略](cluster-container-schedule.md)