# Private Registry

&ensp;&ensp;&ensp;The private registry provides the image storage infrastructure for Humpback, Humpback involves image searching, container building, and so on, all dependent on private registry services.   

##  Prerequisites   

- Docker version is `1.8.3` or later  

- Docker official registry `Distribution`, recommended` 2.5.1` or higher   

##  Start the Registry

&ensp;&ensp;&ensp;This section focuses on building and running Docker private registries in containers, using the `2.5.1` version as an example.  

&ensp;&ensp;&ensp;The image name of Docker official registry `Distribution` at hub.docker.com is <a href="https://hub.docker.com/r/library/registry/">`registry` </a>

&ensp;&ensp;&ensp;1、Pull the registry image to the local

```bash
$ docker pull registry:2.5.1
```
&ensp;&ensp;&ensp;2、Start the registry service   

```bash
$ docker run -d -p 5000:5000 --restart=always \
  -v /var/lib/registry/:/var/lib/registry/ \
  -v /etc/docker/registry/config.yml:/etc/docker/registry/config.yml \
  --name registry \
  registry:2.5.1
```
&ensp;&ensp;&ensp;Refer to the registry configuration instructions <a href="https://github.com/docker/distribution/blob/master/docs/configuration.md">`configuration.md` </a>

&ensp;&ensp;&ensp;3、Verify the registry service 

```bash
$ curl http://localhost:5000/v2/_catalog
{"repositories":[]}
```
&ensp;&ensp;&ensp;If you can normally access the `registry` interface `_catalog`, the service starts successfully. 

##  Configure the Docker  

&ensp;&ensp;&ensp;If you want the Docker daemon to `pull` or` push` image from the registry smoothly, you need to modify the Docker configuration.

&ensp;&ensp;&ensp;In the following configuration instructions, `<registry_server>` is the IP address or host domain name of the registry server.  

- CentOS 6   

  Modify /etc/sysconfig/docker   

  In the `other_args` option add "--insecure-registry `<registry_server>`:5000"

  If there is no `other_args` option for Docker version differences, need to find `INSECURE_REGISTRY` option and modified to: INSECURE_REGISTRY="`<registry_server>`:5000"  

  Note: Part of the docker version of `1.10.x` or higher has default to` CentOS 7` installation source, if the system is `CentOS 7` also need to modify the` docker.service` file.   

  In the `CentOS 7` system, first see if there is a` EnvironmentFile` configuration in /usr/lib/systemd/system/docker.service.

&ensp;&ensp;&ensp;&ensp;&ensp;![C6配置](./_media/centos7_docker_etc.jpg)   
  
&ensp;&ensp;&ensp;&ensp;&ensp;File content may be slightly different, just pay attention to the `EnvironmentFile` and` ExecStart` part in the red box, make sure to link to the configuration file path and the relevant environment variable `INSECURE_REGISTRY`.

- Ubuntu

  Modify /etc/default/docker

  Find the `INSECURE_REGISTRY` option and change to: INSECURE_REGISTRY =" `<registry_server>`: 5000 "

&ensp;&ensp;&ensp;&ensp;After the configuration, restart the Docker local service and check whether the configuration takes effect. If the process appears, the `--insecure-registry` parameter indicates that the configuration takes effect.

```bash
$ service docker restart
$ ps aux | grep docker
root   5003  1.8  2.0  520284  42360 ?  Ssl  15:59  0:00  /usr/bin/dockerd -H fd:// --insecure-registry 192.168.1.10:5000
```

##  Access to Humpback

![C6配置](./_media/system_config.png)    

&ensp;&ensp;&ensp;Log in to Humpback as an administrator and expand the left `Manage` and click `System Config`. Enter the system configuration interface, check the `Enable Private Registry`, and fill in the private registry service address and `Save`.
