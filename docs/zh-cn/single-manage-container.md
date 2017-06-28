# 容器管理

> 容器列表

![Container List](_media/single-container-list.png)

通过列表可以查看当前服务器上运行的所有容器，可以对其进行类似启停、删除的操作，也可以快捷跳转到容器监控和日志查看页面    
*根据屏幕分辨率不同，列表展示的数据可能会有些不同*

> 容器详细信息

![Container Detail](_media/single-container-detail.png)

相比`容器列表`，这里会展现出更详细容器的信息，包括Pid、端口、卷标等内容，除了基本的（Start/Stop/Restart/Pause/UnPause/Delete）等功能外，还有`Rename`、`Upgrade`等功能

- `Rename`：修改容器名称
- `Upgrade`：升/降级容器所使用的镜像，可以对当前组内的服务器进行批量操作，前提是容器名称需要一致

> 容器监控

![Container Monitor](_media/container-monitor.png)

提供对容器CPU和内存的监控

> 容器日志

![Container Logs](_media/single-container-logs.png)

可以根据Tail值确定显示多少条
