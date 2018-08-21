# Cluster And Scheduling   

![Cluster mode](./_media/cluster-mode.png)

##  Cluster mode

The Group can be set to `ClusterMode` . Once set to cluster mode, all the server under the group will be managed and scheduled by  humpback center.

All managed nodes are registered in the cluster via Humpback Agent and zookeeper. The same node can be registered to multiple groups.

Humpback Center is only a scheduler, not responsible for container management.

## How to start the cluster center service   

Run Humpback Center in docker   
```bash 
$ docker pull humpbacks/humpback-center:latest
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-center \
 -e HUMPBACK_SITEAPI=http://192.168.2.80/api \
 -e CENTER_LISTEN_PORT=:8589 \
 -e CENTER_API_ENABLECORS=true \
 -e DOCKER_CLUSTER_URIS=etcd://192.168.2.80:2379,192.168.2.81:2379,192.168.2.82:2379 \
 -e DOCKER_CLUSTER_NAME=humpback/center \
 -v /opt/app/humpback-center/cache:/opt/humpback-center/cache \
 -v /opt/app/humpback-center/logs:/opt/humpback-center/logs \
 -v /opt/app/humpback-center/data:/opt/humpback-center/data \
 -v /opt/app/humpback-center/etc/config.yaml:/opt/humpback-center/etc/config.yaml \
 humpbacks/humpback-center:latest
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND                  CREATED         STATUS         PORTS         NAMES
a1640bf8c956    humpbacks/humpback-center:latest  "./humpback-center"     15 minutes ago  45 seconds ago              humpback-center
```   

## How to join the cluster

Run Humpback Agent in docker
```bash 
$ docker pull humpbacks/humpback-agent:latest
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-agent \
 -e DOCKER_API_VERSION=v1.21 \
 -e DOCKER_AGENT_IPADDR=0.0.0.0 \
 -e DOCKER_CLUSTER_ENABLED=true \
 -e DOCKER_CLUSTER_URIS=etcd://192.168.2.80:2379,192.168.2.81:2379,192.168.2.82:2379 \
 -e DOCKER_CLUSTER_NAME=humpback/center \
 -v /var/run/:/var/run/:rw \
 humpbacks/humpback-agent:latest
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND               CREATED        STATUS         PORTS         NAMES
b1ac4a82c2dd    humpbacks/humpback-agent:latest   "./humpback-agent"   3 minutes ago  20 seconds ago               humpback-agent
```
Related environment variables   

   - HUMPBACK_SITEAPI=http://192.168.2.80/api， Humpback-Web site address, pay attention to bring `/ api`.   

   - CENTER_LISTEN_PORT=:8589, The default port for the Humpback Center API is 8589.   

   - CENTER_API_ENABLECORS=true, Whether the Humpback Center API supports cross-domain access.   

   - DOCKER_API_VERSION=v1.21, Default `v1.21`, use the` docker version` command to view the API Version in the client first.   

   - DOCKER_CLUSTER_ENABLED=true, Cluster mode switch, After the container is started, it is automatically registered in the cluster   

   - DOCKER_CLUSTER_URIS, Default to etcd cluster address, if it is zookeeper, consul or other found tools, please change the prefix to `zk: //` or `consul: //`.   

   - DOCKER_CLUSTER_NAME=humpback/center, Set the humpback cluster name, and the configuration of `Humpback Agent` is consistent with the` Humpback Center` configuration

   - /opt/app/humpback-center/cache Cluster container information persistence directory, It is recommended that you do not manually change and delete.   

   - /opt/app/humpback-center/logs System log directory.     

   - /opt/app/humpback-center/data System data directory.   

After successful startup of Humpback Agent, you can modify the Group attribute at the Humpback site and open the ClusterMode mode option. At this time all the servers under the Group will switch to cluster scheduling mode.

## How to exit the cluster     

When the ClusterMode option of Group is closed on the Humpback site, all services in the group will exit the cluster management mode and do not need to restart the Humpback Agent. Note: Once group exit cluster, all containers which are scheduled by the cluster will be deleted.

## Scheduling strategy   

At present, ClusterMode has only one built-in scheduling strategy. The strategy compute rank according to node's CPU, memory usage, and number of containers.

![Scheduling strategy](./_media/humpback-schduler.png)   

Scheduling a container will use following three steps: 

- `Health Check`：First choose the health nodes.

- `Weight Assessment`：Compute the rank.

- `Node Filtering`：Ignore the nodes which were selected before or create container failed.


## About the container   

#### Container migration   

Container migration occurs when the node is down or Humpback Agent stopped. The scheduler will create new containers in the health nodes. It will make sure the number of running containers is same as you configured. But the migration action will not start as soon as node down, it has a delay time(default 45 seconds). Delay time is used for avoid network unstable.

#### Container recovery   

Humpback Center will scan the number of container instances for all groups every 120s, in order to check the number of container instances. This situation generally occurs when the container is deleted manually. The container recovery algorithm is also same as scheduling strategy.

## About WebHook   

The WebHook is used to call back to notify the cluster container status changed. You can get the detail infomation about the IP address and name of the target host to be assigned.

WebHook event description:
```
CreateMetaEvent    Triggers when the Containers are created for the first time
RemoveMetaEvent    Triggered when removing Containers, including exit clustering mode
OperateMetaEvent   Triggered when the Containers are manipulated, Start/Stop/Restart/Kill…
UpdateMetaEvent    Triggered when modifying the number of cluster instances
UpgradeMetaEvent   Triggered when upgrading Containers Tag
MigrateMetaEvent   Triggered when a container migration occurs
RecoveryMetaEvent  Triggered when the system recovers or shrinks the container
```
WebHook content format:   

The following is the WebHook event for the container to migrate. The container CLUSTER-dc0d7d79-cluster-test-2 was migrated to the 192.168.0.1 node.
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
