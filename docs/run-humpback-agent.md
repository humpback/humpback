# Deploy the Humpback Agent 

> Humpback Agent 

Three servers: 192.168.2.80, 192.168.2.81, 192.168.2.82 need to deploy Humpback Agent for the local image and container management.

- Check the Docker version

First check the `Docker API` version number:` 1.21`. Different Docker versions has different api.

```bash
$ docker version
Client:
 Version:      1.9.1
 API version:  1.21
 Go version:   go1.4.3
 Git commit:   a34a1d5
 Built:        Fri Nov 20 17:56:04 UTC 2015
 OS/Arch:      linux/amd64

Server:
 Version:      1.9.1
 API version:  1.21
 Go version:   go1.4.3
 Git commit:   a34a1d5
 Built:        Fri Nov 20 17:56:04 UTC 2015
 OS/Arch:      linux/amd64
```

- Start the Humpback Agent

 1、Start in cluster mode 
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
 2、Start in single mode
```bash 
$ docker pull humpbacks/humpback-agent:1.0.0
$ docker run -d -ti --net=host --restart=always \
 --name=humpback-agent \
 -e DOCKER_API_VERSION=v1.21 \
 -v /var/run/:/var/run/:rw \
 humpbacks/humpback-agent:1.0.0
$ docker ps -a
CONTAINER ID    IMAGE                           COMMAND               CREATED        STATUS         PORTS         NAMES
b1ac4a82c2dd    humpbacks/humpback-agent:1.0.0   "./humpback-agent"   3 minutes ago  20 seconds ago               humpback-agent
```

- Environment variables

`DOCKER_API_VERSION=v1.21` must be consistent with the above version number.   

`DOCKER_CLUSTER_ENABLED=true`If the current agent needs to be scheduled by cluster mode, set the cluster mode to true, otherwise, close the node of cluster scheduling to `Single Mode` mode.    

`DOCKER_CLUSTER_URIS=zk://192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181` is the previously configured` Zookeeper` cluster address information.   

`DOCKER_CLUSTER_NAME=humpback/center` Cluster name, It is consistent with the `Humpback Center` configuration.

- Register the server and create group(cluster)

After the three servers have successfully started the `Humpback Agent` container, go to the `Humpback Web` site. Expand the left side of the `Manage` function, click `Groups` to create a group `MyCluster`   

And add three node servers to the group and click `Save`.

![Humpback Add Group](_media/humpbackadd-group.png)

Go to the grouping page to view cluster server information.

![Humpback Group](_media/humpback-group.png)

- [How to create a container](single-create-container.md)

- [How to manage containers](single-manage-container.md)

- [Batch operation](single-batch-operate.md)
