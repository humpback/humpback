# Deploy Humpback Center

> Humpback Center  

&ensp;&ensp;&ensp;Humpback Center mainly provides the cluster container scheduling service for the Humpback platform, and as a cluster center role to achieve the various groups of container distribution management.

&ensp;&ensp;&ensp;Need cluster scheduling `Group`, you can open the` Cluster Mode` switch in the settings.

- Start the Humpback Center

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

- Environment variables and parameters

&ensp;&ensp;&ensp;`HUMPBACK_SITEAPI=http://192.168.2.80:8012/api` Humpback-Web site address, pay attention to bring `/ api`.   

&ensp;&ensp;&ensp;`CENTER_LISTEN_PORT=:8589` The default port for the Humpback Center API is: 8589.   

&ensp;&ensp;&ensp;`CENTER_API_ENABLECORS=true`Whether the Humpback Center API supports cross-domain access.

&ensp;&ensp;&ensp;`DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181` is the previously configured `Zookeeper` cluster address information.

&ensp;&ensp;&ensp;`DOCKER_CLUSTER_NAME=humpback/center` The cluster name. It is consistent with the `Humpback Agent` configuration.   

&ensp;&ensp;&ensp;`-v /opt/app/humpback-center/cache` Cluster container information persistence directory. It is recommended that you do not manually change and delete.   

&ensp;&ensp;&ensp;`-v /opt/app/humpback-center/logs` System log directory.

- [How to create a container in cluster mode](cluster-create-container.md)

- [Cluster scheduling strategy](cluster-container-schedule.md)
