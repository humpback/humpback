# How To Start Humpback

> Preparation

First prepare the server. We will use the following three servers to finish our example. Please install the `Docker` for each machine before start. 

```
SERVER01：192.168.2.80 
SERVER02：192.168.2.81
SERVER03：192.168.2.82 
```

> Deployment introduction

`192.168.2.80` as the Humpback platform center will start these service: 

Docker Registry. Image private registry service.    

Humpback-Web. Humpback manage the site.      

Humpback-Center. Provides container scheduling services for cluster mode.   

All servers start `Humpback-Agent` service.   

```
| Server           |  Zookeeper |  Docker Registry  |  Humpback-Web  |  Humpback-Agent  |  Humpback-Center  |
|------------------|:----------:|:-----------------:|:--------------:|:----------------:|:-----------------:|
|   192.168.2.80   |     √      |         √         |        √       |        √         |         √         |
|   192.168.2.81   |     √      |         X         |        X       |        √         |         X         |
|   192.168.2.82   |     √      |         X         |        X       |        √         |         X         |
```

There are two modes in Humpback platform to manage containers:  

`Single Mode`, In this mode, `Humpback-Web` will call the` Humpback-Agent` API directly for simple cluster container management. this mode is simple for deployment, `Humpback-Center` and ` Zookeeper` are not necessary.

`Cluster Mode`, `Humpback-Center` and` Zookeeper` are need for supporting cluster. `Humpback-Center` is a cluster scheduling center, and` Zookeeper` is used for node discovery.

> Zookeeper cluster

[How to build a Zookeeper cluster](run-zookeeper.md) 

> Docker Registry private registry

[How to deploy Docker Registry](run-registry.md)

> Humpback Web site

[How to deploy the Humpback site](run-humpback-web.md)

> Humpback Agent

[How to deploy the Humpback Agent](run-humpback-agent.md)

> Humpback Center

[How to deploy Humpback Center](run-humpback-center.md)









 
