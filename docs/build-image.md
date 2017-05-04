# 创建镜像

&ensp;&ensp;&ensp;在 Humpback 系统中我们虽然可以使用一些公网镜像，但是总有不尽人意的地方，比如镜像拉取耗时或容器环境差异化配置等，那我们就需要自己 build 镜像，并提交到私有镜像仓库中保存，目前有2种方式可以创建镜像。   

- 可以更新现有的已经存在的镜像   

- 通过 `Dockfile` 来创建一个全新的镜像

## 更新现有镜像创建提交一个新镜像   

&ensp;&ensp;&ensp;该方式实质是修改现有正在运行的 `Contaienr` 的  继承目前运行的容器并固化为一个 首选我们先让自己你要更新的镜像跑起来

## 通过 `Dockfile` 制作镜像   


参见官方 <a href="https://docs.docker.com/engine/reference/builder/">`Dockerfile 手册`</a>
https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/
http://qiita.com/xiangzhuyuan/items/0c51fbcd4488095e0687


