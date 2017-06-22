# 部署私有仓库

> Docker Registry

- 拉取 Docker 官方仓库镜像   

&ensp;&ensp;&ensp;以下选用 `192.168.2.80` 作为 Humpback 私有仓库服务器。

```bash
192.168.2.80
$ mkdir -p /var/lib/registry
$ docker pull registry:2.5.1
```

- 启动仓库服务

&ensp;&ensp;&ensp;请将配置文件拷贝到 /etc/docker/registry/config.yml，关于仓库配置说明请参见 <a href="https://github.com/docker/distribution/blob/master/docs/configuration.md">`configuration.md` </a>

&ensp;&ensp;&ensp;仓库若要支持跨域访问，请将 `config.yml` 中 `http` 项改为如下：   

```bash
http:
  addr: :5000
  secret: abcdfg
  headers:
    X-Content-Type-Options: [nosniff]
    Access-Control-Allow-Headers: ['Origin,Accept,Content-Type,Authorization']
    Access-Control-Allow-Origin: ['*']
    Access-Control-Allow-Methods: ['GET,POST,PUT,DELETE']
```
&ensp;&ensp;&ensp;启动仓库   

```bash
$ docker run -d -p 5000:5000 --restart=always \
 -v /var/lib/registry/:/var/lib/registry/ \
 -v /etc/docker/registry/config.yml:/etc/docker/registry/config.yml \
 --name registry \
 registry:2.5.1
```
- 检查启动状态，若 `_catalog` 接口能正常访问，证明仓库启动成功。    

```bash
$ docker ps -a
CONTAINER ID    IMAGE            COMMAND                  CREATED         STATUS            PORTS                    NAMES
a0e4c87eef5b    registry:2.5.1   "/entrypoint.sh /etc/"   15 minutes ago  45 seconds ago    0.0.0.0:5000->5000/tcp   registry
$ curl http://192.168.2.80:5000/v2/_catalog
{"repositories":[]}
``` 