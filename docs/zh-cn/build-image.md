# 创建镜像

&ensp;&ensp;&ensp;在 Humpback 系统中我们虽然可以使用一些公网镜像，但是总有不尽人意的地方，比如镜像拉取耗时或容器配置差异化等因素，此时我们就需要自己定制镜像，并提交到私有镜像仓库中保存，目前有2种方式可以创建镜像。   

- 可以更新现有容器创建提交一个新镜像   

- 通过 `Dockfile` 来创建一个全新的镜像

## 更新现有容器创建提交一个新镜像   

&ensp;&ensp;&ensp;该方式实质是调整现有正在运行的 `Container`，修改其中的配置或相关运行环境，然后固化为一个新版本的 `Image` 到本地。   

&ensp;&ensp;&ensp;首先我们将要更新的 `Image` 成功运行起来成为一个 `Container`，以官方基础镜像 `alpine:3.5` 为例：

```bash
$ docker run -t -i --name=alpine alpine:3.5 /bin/ash
```
&ensp;&ensp;&ensp;然后我们就可以修改这个容器了，比如给容器安装一个 `postgresql-dev` 环境，目的是制作一个基于 `alpine` 的 `postgresql` 基础镜像:   
```bash
/ # apk add --no-cache postgresql-dev
fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/community/x86_64/APKINDEX.tar.gz
(1/9) Installing libressl2.4-libtls (2.4.4-r0)
(2/9) Installing pkgconf (1.0.2-r0)
(3/9) Installing libressl-dev (2.4.4-r0)
(4/9) Installing db (5.3.28-r0)
(5/9) Installing libsasl (2.1.26-r8)
(6/9) Installing libldap (2.4.44-r3)
(7/9) Installing libpq (9.6.2-r0)
(8/9) Installing postgresql-libs (9.6.2-r0)
(9/9) Installing postgresql-dev (9.6.2-r0)
Executing busybox-1.25.1-r0.trigger
OK: 29 MiB in 20 packages
```   
&ensp;&ensp;&ensp;执行 `exit` 退出这个容器，然后可以通过 `docker commit` 命令来将当前 `alpine` 容器固化为一个新的镜像了。   

&ensp;&ensp;&ensp;退出容器后如下所示，查看被更新后的容器编号为 `b1ac4a82c2dd`
```bash
/ # exit
$ docker ps -a
CONTAINER ID    IMAGE         COMMAND       CREATED          STATUS                   PORTS      NAMES
b1ac4a82c2dd    alpine:3.5    "/bin/ash"    "5 minutes ago"  Exit (0) 5 seconds ago              alpine
```
&ensp;&ensp;&ensp;执行 `docker commit` 命令将容器 `b1ac4a82c2dd` 固化为新镜像:   

```bash
$ docker commit -m "add postgresql" -a "bobliu" b1ac4a82c2dd 192.168.1.10:5000/postgresql:9.6.2
sha256:b1ac4a82c2dd17e9c50b6bff0e42290496df9b6187117dfda9daec98015f6ae1
$ docker images
REPOSITORY                      TAG           IMAGE ID          CREATED          SIZE
alpine                          3.5           1bb3a95866d7      15 minutes ago   3.987MB
192.168.1.10:5000/postgresql    9.6.2         b1ac4a82c2dd      45 seconds ago   27.59MB
```
&ensp;&ensp;&ensp;注意 `docker commit` 命令格式如下：
```bash
docker commit [OPTIONS] CONTAINER [REPOSITORY[:TAG]]
```
&ensp;&ensp;&ensp;主要选项(OPTIONS)如下：   
- -a, --author - {string}, 作者（如："John Hannibal Smith "）
- -c, --change - {list}, 使用Dockerfile指令来创建镜像（默认 []）
- -m, --message - {string}, 提交备注信息
- -p, --pause - {string}, 提交时暂停容器（默认 true）

&ensp;&ensp;&ensp;紧接 `CONTAINER` 是被固化的容器ID `b1ac4a82c2dd`。   

&ensp;&ensp;&ensp;最后 `[REPOSITORY[:TAG]]` 是新镜像名称与 tag: `9.6.2`，示例中`192.168.1.10:5000` 为私有仓库地址。

```bash
192.168.1.10:5000/postgresql:9.6.2
```

&ensp;&ensp;&ensp;最后将镜像提交到私有仓库 `192.168.1.10:5000`   

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6.2
The push refers to a repository [192.168.1.10:5000/postgresql]
c1bfc2be7117: Pushed
23b9c7b43573: Pushed
9.6.2: digest: sha256:afcf3f596e42837ded618f4694e6ff9928034ce20e860b75125992c3dc1ba501 size: 739
```

## 通过 `Dockfile` 制作镜像   

&ensp;&ensp;&ensp;使用 `docker commit` 命令很容易基于现有的容器创建新的扩展，除此之外我们还可以通过 `docker build` 来从零开始构建一个镜像，特别是基于一些基础镜像来构建我们自己的应用镜像，使用 `Dockerfile` 和 `docker build` 命令来构建镜像操作更灵活、过程可重复，因此也更推荐使用这种方式来构建镜像。   

&ensp;&ensp;&ensp; `Dockerfile` 基于 `DSL (Domain Specific Language)` 语言构建 `Docker` 镜像，`Dockerfile` 编写完成后，就可以使用 `docker build` 命令来构建一个新镜像。

&ensp;&ensp;&ensp;首先创建一个文件夹，取名为 `postgresql`，这个目录就是我们的构建环境，在 Docker 中，将这个环境称为上下文（content）或者构建上下文（build content），构建镜像时 Docker 会将构建环境中的文件和目录传递给守护进程，这样守护进程就访问到用户想在镜像中存储的任何代码、文件或其它数据。   

&ensp;&ensp;&ensp;在目录内创建一个 `Dockerfile` 文件：

```bash
$ mkdir postgresql
$ cd postgresql
$ touch Dockerfile
```

&ensp;&ensp;&ensp;然后编辑这个 `Dockerfile` 文件，同样以官方 `alpine:3.5` 基础镜像。     

```bash
# This is a comment
FROM alpine:3.5
MAINTAINER bobliu <bobliu@example.com>
RUN apk add --no-cache postgresql-dev
```

&ensp;&ensp;&ensp;如何使用 `Dockerfile` 内部相关命令，请参考官方：<a href="https://docs.docker.com/engine/reference/builder/">`Dockerfile Reference`</a>    

&ensp;&ensp;&ensp;如何书写 `Dockerfile` 请参考官方：<a href="https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/">`Best practices for writing Dockerfiles`</a>   

&ensp;&ensp;&ensp;接下来通过 `docker build` 命令来构建 `postgresql` 镜像:    

```bash
$ docker build -t 192.168.1.10:5000/postgresql:9.6 ./postgresql
Sending build context to Docker daemon 2.048 kB
Step 1 : FROM alpine:3.5
3.5: Pulling from library/alpine

031b4db7df57: Pull complete
Digest: sha256:4d1622c09e7e7128132f825351d61f0e66651512672820a51ba537f0fd673ffb
Status: Downloaded newer image for alpine:3.5
 ---> 031b4db7df57
Step 2 : MAINTAINER bobliu <bobliu@example.com>
 ---> Running in cd85fee20da5
 ---> 51c8ce3d6bb8
Removing intermediate container cd85fee20da5
Step 3 : RUN apk add --no-cache postgresql-dev
 ---> Running in 2948f616151f
fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/community/x86_64/APKINDEX.tar.gz
(1/9) Installing libressl2.4-libtls (2.4.4-r0)
(2/9) Installing pkgconf (1.0.2-r0)
(3/9) Installing libressl-dev (2.4.4-r0)
(4/9) Installing db (5.3.28-r0)
(5/9) Installing libsasl (2.1.26-r8)
(6/9) Installing libldap (2.4.44-r3)
(7/9) Installing libpq (9.6.2-r0)
(8/9) Installing postgresql-libs (9.6.2-r0)
(9/9) Installing postgresql-dev (9.6.2-r0)
Executing busybox-1.25.1-r0.trigger
OK: 29 MiB in 20 packages
 ---> c3e27ed7f727
Removing intermediate container 2948f616151f
Successfully built c3e27ed7f727
```  

&ensp;&ensp;&ensp;构建成功后，最后将镜像提交到私有仓库 `192.168.1.10:5000`   

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6
The push refers to a repository [192.168.1.10:5000/postgresql]
031b4db7df57: Pushed
51c8ce3d6bb8: Pushed
9.6: digest: sha256:c3e27ed7f727d7b0248894ddfa6da54d913abc79e49467ff7e311c1dcd23ffd0 size: 739
```








