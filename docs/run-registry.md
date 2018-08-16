# Deploy a private registry

> Docker Registry

- Pull the Docker official registry image   

We use `192.168.2.80` as Humpback private registry server.

```bash
192.168.2.80
$ mkdir -p /var/lib/registry
$ docker pull registry:2.5.1
```

- Start the registry service

Please copy the configuration file to /etc/docker/registry/config.yml. For registry configuration instructions, see <a href="https://github.com/docker/distribution/blob/master/docs/configuration.md">`configuration.md` </a>

If registry need to support CORS, change the `http` entry in` config.yml` as follows:  

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
Start the registry 

```bash
$ docker run -d -p 5000:5000 --restart=always \
 -v /var/lib/registry/:/var/lib/registry/ \
 -v /etc/docker/registry/config.yml:/etc/docker/registry/config.yml \
 --name registry \
 registry:2.6.2
```
- Check the registry state. Check the `_catalog` API as below:  

```bash
$ docker ps -a
CONTAINER ID    IMAGE            COMMAND                  CREATED         STATUS            PORTS                    NAMES
a0e4c87eef5b    registry:2.6.2   "/entrypoint.sh /etc/"   15 minutes ago  45 seconds ago    0.0.0.0:5000->5000/tcp   registry
$ curl http://192.168.2.80:5000/v2/_catalog
{"repositories":[]}
``` 
