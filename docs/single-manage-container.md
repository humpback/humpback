# Container Management

> Container list

![Container List](_media/single-container-list.png)

The list show all the containers running on the current server. You can do some operations such as starting, stopping and deleting it, you can quickly jump to the container monitoring and log viewing page
*Depending on the screen resolution, the data displayed by the list may be slightly different*

> Container details

![Container Detail](_media/single-container-detail.png)

There are more detailed container information than `container list`, including Pid, port, volume, log config and other content. You can do `Rename`, `Upgrade`  and other actions here. 

- `Rename`：Modify the container name
- `Upgrade`：Replace current container with a new one using same images but different tag. The operation can be batch to other servers in current group.

> Container monitoring

![Container Monitor](_media/container-monitor.png)

Provides monitoring of container CPU and memory

> Container log

![Container Logs](_media/single-container-logs.png)

You can determine how many bars to display based on the Tail value
