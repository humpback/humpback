# Create a image

In the Humpback system, although we can use some public network images, but there are always something unsatisfactory. Such as long time when pulling images or container configuration issue. So we need to customize the image, and submitted to the private image registry. There are two ways to create a image.  

- You can create a new image base on an existing container.

- Create a new image with `Dockfile`.

## Create a new image base an existing container  

This way is that adjust an exists running `Container` and modify the configuration or environment args, then create a new image base on it.   

First we start a `Container` from the `Image` which need to be updated. Use the official base mirror `alpine: 3.5` as a example.

```bash
$ docker run -t -i --name=alpine alpine:3.5 /bin/ash
```
And then we can modify the container, for example instal a `postgresql-dev` environment for the container:

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
Exit this container and then use the `docker commit` command to crate a new image base on the current `alpine` container.   

```bash
/ # exit
$ docker ps -a
CONTAINER ID    IMAGE         COMMAND       CREATED          STATUS                   PORTS      NAMES
b1ac4a82c2dd    alpine:3.5    "/bin/ash"    "5 minutes ago"  Exit (0) 5 seconds ago              alpine
```
Execute the `docker commit` command to create new image:   

```bash
$ docker commit -m "add postgresql" -a "bobliu" b1ac4a82c2dd 192.168.1.10:5000/postgresql:9.6.2
sha256:b1ac4a82c2dd17e9c50b6bff0e42290496df9b6187117dfda9daec98015f6ae1
$ docker images
REPOSITORY                      TAG           IMAGE ID          CREATED          SIZE
alpine                          3.5           1bb3a95866d7      15 minutes ago   3.987MB
192.168.1.10:5000/postgresql    9.6.2         b1ac4a82c2dd      45 seconds ago   27.59MB
```
The `docker commit` command has the following options:
```bash
docker commit [OPTIONS] CONTAINER [REPOSITORY[:TAG]]
```
The main options (OPTIONS) are as below:   
- -a, --author - {string}, The author
- -c, --change - {list}, Use the Dockerfile directive to create a image (default [])
- -m, --message - {string}, Submit a note message
- -p, --pause - {string}, Suspended container when submitting (default true)

The example` 192.168.1.10: 5000` is the private registry address.

```bash
192.168.1.10:5000/postgresql:9.6.2
```

Finally, push image to the private registry `192.168.1.10: 5000` 

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6.2
The push refers to a repository [192.168.1.10:5000/postgresql]
c1bfc2be7117: Pushed
23b9c7b43573: Pushed
9.6.2: digest: sha256:afcf3f596e42837ded618f4694e6ff9928034ce20e860b75125992c3dc1ba501 size: 739
```

## Create a new image with `Dockfile`

`Dockerfile` builds the `Docker` image based on the `DSL (Domain Specific Language)` language. After the `Dockerfile` is written, you can use the `docker build` command to build a new image.

First create a folder, named `postgresql`, this directory is our build environment, in the Docker, this environment will be called content or build content. When building a image, the Docker passes the files and directories in the build environment to the daemons so that the daemon accesses any code, files, or other data that the user wants to store in the image.   

Create a `Dockerfile` file in the directory:

```bash
$ mkdir postgresql
$ cd postgresql
$ touch Dockerfile
```

Then edit the `Dockerfile` file, again, with the official `alpine: 3.5` as the base image.  

```bash
# This is a comment
FROM alpine:3.5
MAINTAINER bobliu <bobliu@example.com>
RUN apk add --no-cache postgresql-dev
```

How to use `Dockerfile` commands, please refer to the official:<a href="https://docs.docker.com/engine/reference/builder/">`Dockerfile Reference`</a>    

How to write `Dockerfile`, please refer to the official:<a href="https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/">`Best practices for writing Dockerfiles`</a>   

Then use the `docker build` command to build `postgresql` image:    

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

Then you can push the image to the private registry `192.168.1.10: 5000`  

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6
The push refers to a repository [192.168.1.10:5000/postgresql]
031b4db7df57: Pushed
51c8ce3d6bb8: Pushed
9.6: digest: sha256:c3e27ed7f727d7b0248894ddfa6da54d913abc79e49467ff7e311c1dcd23ffd0 size: 739
```








