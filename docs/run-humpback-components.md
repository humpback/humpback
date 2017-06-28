# How To Start Humpback

> Ready to work   

First prepare the server, the following three machines as a build cluster example. Please install the `Docker` instance for each machine before you start. 

```
SERVER01：192.168.2.80 
SERVER02：192.168.2.81
SERVER03：192.168.2.82 
```

> Deployment introduction

&ensp;&ensp;&ensp;You can see `192.168.2.80` as the Humpback platform center, the largest number of services to start: 

&ensp;&ensp;&ensp;Docker Registry. Image private registry service.    

&ensp;&ensp;&ensp;Humpback-Web. Humpback manage the site.      

&ensp;&ensp;&ensp;Humpback-Center. Provides container scheduling services for cluster mode.   

&ensp;&ensp;&ensp;At the same time the following three servers are activated `Humpback-Agent` to provide local image management function, waiting for the Humpback platform call.   

```
| Server           |  Zookeeper |  Docker Registry  |  Humpback-Web  |  Humpback-Agent  |  Humpback-Center  |
|------------------|:----------:|:-----------------:|:--------------:|:----------------:|:-----------------:|
|   192.168.2.80   |     √      |         √         |        √       |        √         |         √         |
|   192.168.2.81   |     √      |         X         |        X       |        √         |         X         |
|   192.168.2.82   |     √      |         X         |        X       |        √         |         X         |
```

&ensp;&ensp;&ensp;Each server needs to start the `Zookeeper` node to form a cluster, but building a` Zookeeper` cluster is not a must. 
  
&ensp;&ensp;&ensp;Humpback platform management mode is divided into two types:  

&ensp;&ensp;&ensp;`Single Mode`, In this mode, `Humpback-Web` will directly access the` Humpback-Agent` API for interactive implementation of simple cluster container management, this mode is deployed more single, do not have to start `Humpback-Center` and build` Zookeeper` cluster.

&ensp;&ensp;&ensp;`Cluster Mode`, Humpback back-end needs to start `Humpback-Center` and` Zookeeper` cluster to support, `Humpback-Center` as a cluster dispatch center to provide container dynamic scheduling allocation function, and` Zookeeper` cluster provides node discovery function, to achieve cluster node flexibility telescopic.

> Zookeeper cluster

&ensp;&ensp;&ensp;[How to build a Zookeeper cluster](run-zookeeper.md) 

> Docker Registry private registry

&ensp;&ensp;&ensp;[How to deploy Docker Registry](run-registry.md)

> Humpback Web site

&ensp;&ensp;&ensp;[How to deploy the Humpback site](run-humpback-web.md)

> Humpback Agent

&ensp;&ensp;&ensp;[How to deploy the Humpback Agent](run-humpback-agent.md)

> Humpback Center

&ensp;&ensp;&ensp;[How to deploy Humpback Center](run-humpback-center.md)









 
