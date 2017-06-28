# Humpback Architecture

![Architecture diagram](./_media/humpback-arch.png)

## System role

- `Humpback Website`  
   Humpback management site mainly provides system visual management, it has group management, rights management, registry image query, cluster management functions.

- `Humpback Center`  
   Humpback cluster center, which through the scheduling strategy to implement the container volume creation, operation, delete, upgrade, fault migration, instance adjustment and other functions for cluster; also responsible for cluster node discovery and management.  

- `Humpback Agent`   
   Humpback cluster nodes, multiple nodes in the cluster waiting to be scheduled and managed in the `Humpback Center` group, registered to the cluster center through the node discovery module and maintain heartbeat.    
   
   When the packet is managed, the IP address corresponding to the `Humpback Agent` node can be added to multiple packets. When the container is scheduled, the containers between groups and groups are not isolated from each other.

## Toolset  
   
- `Zookeeper`   

   Humpback defaults to `Zookeeper` for cluster node registration and discovery, and the backend needs to have a `Zookeeper` cluster.   
   
   Default version `3.4.6`, download <a href="http://apache.org/dist/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz">`zookeeper-3.4.6.tar.gz`</a>
   
- `Docker Images Repository`   
 
  Humpack uses `Docker` official private image registry `distribution` to provide image storage services, `Humpback Website` provides access to the registry image query by visiting the `distribution` API.   
      
  Default version of `2.5.1` or higher, see more <a href="https://github.com/docker/distribution/blob/master/README.md">`Docker Distribution`</a>
