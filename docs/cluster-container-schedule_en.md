# Cluster And Scheduling   

![Cluster mode](./_media/cluster-mode.png)

##  Cluster mode

&ensp;&ensp;&ensp;When using Humpback to manage the group, the Group can be set to `ClusterMode` mode according to the business needs. Once set to cluster mode, all the server nodes Humpback Agent under the group will accept the humpback center service management schedule.

&ensp;&ensp;&ensp;All managed Humpback Agent nodes are registered in the cluster with heartbeat, and then the Humpback Center will discover the modules and nodes according to the node and is built in the internal `EnginesPool` for state management. The same Humpback Agent node can be deployed at the same time through the Humpback site to multiple groups.   

&ensp;&ensp;&ensp;Humpback Center is just a scheduler and does not run the container itself, only accept API requests, through the scheduling strategy to select the appropriate Engine in the group to deal with local Containers. This means that even if the Humpback Center is down or shut down for some reason, the containers that are already running will not have any effect for a short period of time.

## How to start the cluster center service   

&emsp;Run Humpback Center in docker   
```bash 
$ docker pull humpbacks/humpback-center:1.0.0
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-center \
 -e HUMPBACK_SITEAPI=http://192.168.2.80/api \
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

## How to join the cluster

&emsp;Run Humpback Agent in docker
```bash 
$ docker pull humpbacks/humpback-agent:1.0.0
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-agent \
 -e DOCKER_API_VERSION=v1.21 \
 -e DOCKER_CLUSTER_ENABLED=true \
 -e DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181 \
 -e DOCKER_CLUSTER_NAME=humpback/center \
 -v /var/run/:/var/run/:rw \
 humpbacks/humpback-agent:1.0.0
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND               CREATED        STATUS         PORTS         NAMES
b1ac4a82c2dd    humpbacks/humpback-agent:1.0.0   "./humpback-agent"   3 minutes ago  20 seconds ago               humpback-agent
```
&ensp;&ensp;&ensp;Related environment variables   

   - HUMPBACK_SITEAPI=http://192.168.2.80/api， Humpback-Web site address, pay attention to bring `/ api`.   

   - CENTER_LISTEN_PORT=:8589, The default port for the Humpback Center API is 8589.   

   - CENTER_API_ENABLECORS=true, Whether the Humpback Center API supports cross-domain access.

   - DOCKER_API_VERSION=v1.20, Default `v1.20`, use the` docker version` command to view the API Version in the client first.

   - DOCKER_CLUSTER_ENABLED=true, Cluster mode switch, After the container is started, it is automatically registered in the cluster   

   - DOCKER_CLUSTER_URIS, Default to zookeeper cluster address, if it is etcd, consul or other found tools, please change the prefix to `etcd: //` or `consul: //`.   

   - DOCKER_CLUSTER_NAME=humpback/center, Set the humpback cluster name, and the configuration of `Humpback Agent` is consistent with the` Humpback Center` configuration

   - /opt/app/humpback-center/cache Cluster container information persistence directory, It is recommended that you do not manually change and delete.   

   - /opt/app/humpback-center/logs System log directory。  

&ensp;&ensp;&ensp;当 Humpback Agent After successful startup, you can modify the Group attribute at the Humpback site and open the ClusterMode mode option. At this point, all the servers under the Group will switch to cluster scheduling mode.

## How to exit the cluster     

&ensp;&ensp;&ensp;在 Humpback 站点将 Group 的 ClusterMode 选项关闭，此时 Group 中所有服务将退出集群管理模式，不需要重启 Humpback Agent 容器即可，若要彻底退出集群避免心跳，可以将 Humpback Agent 容器停止后删除重建，重建时将 `DOCKER_CLUSTER_ENABLED` 环境变量设置为false即可。需要注意：一但对出集群，与该 Group 相关所有之前受集群调度的容器将会被统一删除。   

## Scheduling strategy   

&ensp;&ensp;&ensp;At present, ClusterMode has only one built-in scheduling strategy. The strategy classifies each node weight based on the CPU, Memory usage, and number of Containers per node, and the base information for the node is committed to the cluster at the time of the first start of registration.      

![Scheduling strategy](./_media/humpback-schduler.png)   

&ensp;&ensp;&ensp;When scheduling a container, In the order, the strategy is divided into the following three stages: 

- `Health Check`：Health screening, When the container is dispatched, all valid nodes are selected first, and the principle is the online active node.   

- `Weight Assessment`：Weight assessment, Sort out the available CPU, the largest memory, and the smallest number of containers node.  

- `Node Filtering`：Node filtering, finally, according to Filter filter the node that has been allocated or assigned but failed to schedule.   

&ensp;&ensp;&ensp;If the number of instances (number of instances> cluster servers) can not be assigned after the above three-stage decision, the scheduler randomly selects a node in the node that has been allocated, and in this case, `--net = host` Host mode or specify a fixed port to start the container will be usually scheduled to fail, then enter the Humpback alarm process for e-mail notification.   

&ensp;&ensp;&ensp;At present, only this kind of scheduling strategy algorithm is cured, the purpose is to try to disperse the container in the cluster, the benefits of doing so is that if the node is broken will not lose too much Container, and can play a part of the role of load balancing.

## About the container   

#### Container shrink  

&ensp;&ensp;&ensp;In the case of container shrinkage, this situation usually occurs when the number of instances of the cluster container is modified. When the number of instances is changed from large to small, the scheduling policy selects all the valid nodes every time, compares the nodes with the largest number of containers and arranges them in the reverse order. and eventually dispatches a node with the largest number of local containers to delete the container, the purpose is to try to keep the container dispersed.

#### Container migration   

&ensp;&ensp;&ensp;容器迁移发生在节点宕机或关闭 Humpback Agent 时，目的是按集群容器实例数进行逐步创建恢复，此时调度器会将故障节点上的所有需要被调度的容器进行重新分配，Humpback Center 在检测到节点异常时不会立即触发迁移，而是给出了延时处理（默认为45s），一但超出这个阈值后若节点还未及时恢复上线，此时容器迁移立即触发。这样做的主要目的是为了防止节点抖动导致容器服务间隙中断。

#### Container recovery   

&ensp;&ensp;&ensp;Humpback Center 启动后会定期默认120s对集群所有 Group 进行一次容器实例数扫描，目的是为了将实例数不足的容器恢复出来，这种情况一般出现在分配失败或直接手动删除docker deamon上的容器实例时发生，容器恢复算法也按照调度策略来进行计算。

## About WebHook   

&ensp;&ensp;&ensp;The WebHook function is used to call back to notify the cluster container instance or state to change. At the same time, you can know the IP address and name of the target host to be assigned. When you create a Container, you can set up multiple WebHooks. If you do not set up, the system will abandon the notification after handling the container.

&ensp;&ensp;&ensp;WebHook event description:
```
CreateMetaEvent    Triggers when the Containers are created for the first time
RemoveMetaEvent    Triggered when removing Containers, including exit clustering mode
OperateMetaEvent   Triggered when the Containers are manipulated, Start/Stop/Restart/Kill…
UpdateMetaEvent    Triggered when modifying the number of cluster instances
UpgradeMetaEvent   Triggered when upgrading Containers Tag
MigrateMetaEvent   Triggered when a container migration occurs
RecoveryMetaEvent  Triggered when the system recovers or shrinks the container
```
&ensp;&ensp;&ensp;WebHook content format:   

&ensp;&ensp;&ensp;The following is the WebHook event for the container to migrate. The container CLUSTER-dc0d7d79-cluster-test-2 was migrated to the 192.168.0.1 node.
```js
{
    "Timestamp": 1491976204448307844,
    "Event": "MigrateMetaEvent",
    "MetaBase": {
        "GroupId": "dc0d7d79-e9a1-46aa-8d47-95335f061b60",
        "MetaId": "0136b582-35ca-4755-8148-b7edde2a1127",
        "Instances": 9,
        "WebHooks": [
            {
                "Url": "http://192.168.0.10:8513/mkkks4r/users",
                "SecretToken": "123456"
            }
        ],
        "ImageTag": "latest",
        "Config": {
            "Id": "",
            "Image": "RegistryServer/base/alpine-bash",
            "Command": "date",
            "Name": "cluster-test",
            "Ports": [],
            "Volumes": [],
            "Dns": null,
            "Env": [],
            "HostName": "",
            "NetworkMode": "host",
            "RestartPolicy": "no",
            "Extrahosts": null,
            "Links": [],
            "Ulimits": null
        }
    },
    "HookContainers": [
        {
            "IP": "192.168.0.1",
            "Name": "SERVER01",
            "Container": {
                "Id": "8af07e8091aa34660e7c53b76083a36dd114250548c7df5670e948e7b15eae05",
                "Image": "RegistryServer/base/alpine-bash",
                "Command": "date",
                "Name": "CLUSTER-dc0d7d79-cluster-test-2",
                "Ports": null,
                "Volumes": null,
                "Dns": null,
                "Env": [],
                "HostName": "SCMISMDTSLR02",
                "NetworkMode": "host",
                "Status": {
                    "Status": "exited",
                    "Running": false,
                    "Paused": false,
                    "Restarting": false,
                    "OOMKilled": false,
                    "Dead": false,
                    "Pid": 0,
                    "ExitCode": 0,
                    "Error": "",
                    "StartedAt": "2017-04-12T03:33:31.216194816Z",
                    "FinishedAt": "2017-04-12T03:33:31.230026534Z"
                },
                "RestartPolicy": "no",
                "Extrahosts": null,
                "Links": null,
                "Ulimits": null
            }
        }
    ]
}
```

- `Timestamp`：WebHook triggers timestamps.
- `Event`：Trigger event. See event description.
- `MetaBase`：Meta represents the same batch of containers that contain multiple container instances. In cluster mode MetaID represents a unique batch of containers.
- `GroupID`：is the group to which the container belongs.
- `MetaID`：The container batch ID number is unique in the cluster. If you receive the same MetaID WebHook, you need to determine Timestamp, take the largest for the latest.
- `Instances`：The number of container instances set for this batch.
- `WebHooks`：All WenHook settings are set in the configuration.
- `ImageTag`：The image tag used by the current batch container, ImageTag will change if the same upgrade is made.
- `Config`：creates configuration information for the initial batch.
- `HookContainers`：reports the target node information that the container is assigned: IP and Name, along with container information and status information.   
