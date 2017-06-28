# 管理镜像

> 私有仓库地址设置 [私有仓库搭建](build-registry.md)

![Setting Private Registry](_media/setting-private-registry.png)

> 浏览镜像列表

![Image List](_media/image-list.png)

查看仓库里面的docker镜像，点击`View Detail`可以查看镜像描述、Tag、Dockerfile等信息，其中描述和Dockerfile需要自行维护

> 镜像描述

![Image Description](_media/image-description.png)

使用Markdown维护，主要用于描述镜像用途

> Tag列表

![Image Tags](_media/image-tags.png)

由于Docker官方提供的私有仓库能读取到的镜像信息有限，所以目前只展示了tag

> Dockerfile

![Image Dockerfile](_media/image-dockerfile.png)

维护构建镜像是所使用的Dockerfile，方便使用者查看镜像的构建步骤
