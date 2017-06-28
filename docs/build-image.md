# Create a image

&ensp;&ensp;&ensp;In the Humpback system, although we can use some public network images, but there are always unsatisfactory places, such as taking the image time-consuming long or container configuration differentiation and other factors, then we need to customize the image, and submitted to the private image registry to store, There are currently two ways to create a mirror.  

- You can update an existing container to create a new image   

- Create a brand new image with `Dockfile`

## update an existing container to create a new image   

&ensp;&ensp;&ensp;This way is to adjust the existing running `Container`, modify the configuration or the associated runtime environment, and then solidify it into a new version of `Image` to the local.   

&ensp;&ensp;&ensp;First we run the `Image` to be updated successfully, making it a `Container`, with the official base mirror like 'alpine: 3.5`:

```bash
$ docker run -t -i --name=alpine alpine:3.5 /bin/ash
```
&ensp;&ensp;&ensp;And then we can modify the container, such as install a `postgresql-dev` environment for the container, The purpose is to make a `postgresql` base image based on 'alpine`:

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
&ensp;&ensp;&ensp;Execute `exit` to exit this container and then use the `docker commit` command to solidify the current `alpine` container to a new image.   

&ensp;&ensp;&ensp;After exiting the container, the results are shown below. View the updated container number `b1ac4a82c2dd`
```bash
/ # exit
$ docker ps -a
CONTAINER ID    IMAGE         COMMAND       CREATED          STATUS                   PORTS      NAMES
b1ac4a82c2dd    alpine:3.5    "/bin/ash"    "5 minutes ago"  Exit (0) 5 seconds ago              alpine
```
&ensp;&ensp;&ensp;Execute the `docker commit` command to solidify the container `b1ac4a82c2dd` as a new image:   

```bash
$ docker commit -m "add postgresql" -a "bobliu" b1ac4a82c2dd 192.168.1.10:5000/postgresql:9.6.2
sha256:b1ac4a82c2dd17e9c50b6bff0e42290496df9b6187117dfda9daec98015f6ae1
$ docker images
REPOSITORY                      TAG           IMAGE ID          CREATED          SIZE
alpine                          3.5           1bb3a95866d7      15 minutes ago   3.987MB
192.168.1.10:5000/postgresql    9.6.2         b1ac4a82c2dd      45 seconds ago   27.59MB
```
&ensp;&ensp;&ensp;Note that the `docker commit` command has the following format:
```bash
docker commit [OPTIONS] CONTAINER [REPOSITORY[:TAG]]
```
&ensp;&ensp;&ensp;The main options (OPTIONS) are as follows:   
- -a, --author - {string}, The author（如："John Hannibal Smith "）
- -c, --change - {list}, Use the Dockerfile directive to create a image (default [])
- -m, --message - {string}, Submit a note message
- -p, --pause - {string}, Suspended container when submitting (default true)

&ensp;&ensp;&ensp;Immediately after `CONTAINER` is the solidified container ID `b1ac4a82c2dd`.  

&ensp;&ensp;&ensp;Finally `[REPOSITORY [: TAG]]` is the new image name with tag: `9.6.2`, the example` 192.168.1.10: 5000` is the private registry address.

```bash
192.168.1.10:5000/postgresql:9.6.2
```

&ensp;&ensp;&ensp;Finally, the image will be submitted to the private registry `192.168.1.10: 5000` 

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6.2
The push refers to a repository [192.168.1.10:5000/postgresql]
c1bfc2be7117: Pushed
23b9c7b43573: Pushed
9.6.2: digest: sha256:afcf3f596e42837ded618f4694e6ff9928034ce20e860b75125992c3dc1ba501 size: 739
```

## Make a image with `Dockfile`   

&ensp;&ensp;&ensp;Using the `docker commit` command makes it easy to create new extensions based on existing containers. In addition, we can build a image from scratch at `docker build`, especially based on some basic image to build our own application image, and use the` Dockerfile` and `docker build` commands to build image operations more flexible , The process can be repeated, so it is recommended to use this way to build the image.   

&ensp;&ensp;&ensp;`Dockerfile` builds the `Docker` image based on the `DSL (Domain Specific Language)` language. After the `Dockerfile` is written, you can use the `docker build` command to build a new image.

&ensp;&ensp;&ensp;First create a folder, named `postgresql`, this directory is our build environment, in the Docker, this environment will be called content or build content. When building a image, the Docker passes the files and directories in the build environment to the daemons so that the daemon accesses any code, files, or other data that the user wants to store in the image.   

&ensp;&ensp;&ensp;Create a `Dockerfile` file in the directory:

```bash
$ mkdir postgresql
$ cd postgresql
$ touch Dockerfile
```

&ensp;&ensp;&ensp;Then edit the `Dockerfile` file, again, with the official `alpine: 3.5` as the base image.  

```bash
# This is a comment
FROM alpine:3.5
MAINTAINER bobliu <bobliu@example.com>
RUN apk add --no-cache postgresql-dev
```

&ensp;&ensp;&ensp;How to use `Dockerfile` internal related commands, please refer to the official:<a href="https://docs.docker.com/engine/reference/builder/">`Dockerfile Reference`</a>    

&ensp;&ensp;&ensp;How to write `Dockerfile`, please refer to the official:<a href="https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/">`Best practices for writing Dockerfiles`</a>   

&ensp;&ensp;&ensp;Then through the `docker build` command to build `postgresql` image:    

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

&ensp;&ensp;&ensp;After the build is successful, you will eventually submit the image to the private registry `192.168.1.10: 5000`  

```bash
$ docker push 192.168.1.10:5000/postgresql:9.6
The push refers to a repository [192.168.1.10:5000/postgresql]
031b4db7df57: Pushed
51c8ce3d6bb8: Pushed
9.6: digest: sha256:c3e27ed7f727d7b0248894ddfa6da54d913abc79e49467ff7e311c1dcd23ffd0 size: 739
```








