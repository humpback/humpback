# Humpback Architecture

![Architecture diagram](./_media/humpback-arch.png)

## Components

- `Humpback Website`  
   Humpback management site mainly provides system visual management, it has group management, rights management, registry image query, cluster management functions.

- `Humpback Center`  
   Humpback cluster center, provide container scheduling, batch container operation(create, upgrade, delete), fault migration, instance adjustment etc; also responsible for cluster node discovery and management.  

- `Humpback Agent`   
   A program running in every cluster nodes. Recevie the command from Humpback Web or Humpback Center and execute them. Agent can join any cluster via the service discovery, and a node can join multiple clusters(groups), the containers between each groups are isolated.

## Toolset  
   
- `Zookeeper`   

   Humpback defaults use `Zookeeper` for node registration and discovery.   
   
   Default version `3.4.6`, download <a href="http://apache.org/dist/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz">`zookeeper-3.4.6.tar.gz`</a>
   
- `Docker Images Repository`   
 
  Humpack uses `Docker` official private image registry `distribution` to provide image storage services, `Humpback Website` provides registry image query by access the `distribution` API.   
      
  Default version of `2.5.1` or higher, see more <a href="https://github.com/docker/distribution/blob/master/README.md">`Docker Distribution`</a>
