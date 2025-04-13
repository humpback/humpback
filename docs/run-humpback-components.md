# Install Humpback

## Prerequisites

Humpback consists of two components, the Humpback Server, and the Humpback Agent. Both them run as lightweight Docker containers on a Docker engine.

To get started, you will need the latest version of Docker installed and working. We recommend following the [official installation instructions](https://docs.docker.com/engine/install/) for Docker - in particular, we advise against installing Docker via snap on Ubuntu distributions as you may run into compatibility issues.

## Deployment Hummpback Server

First, create the volume that Humpback Server will use to store its database:

```bash
docker volume create humpback_data
```

Then, install the Portainer Server container:

```bash
docker run -d \
  --name humpback-server \
  -p 8100:8100 \
  -p 8101:8101 \
  --restart=always \
  -v portainer_data:/workspace/data \
  -e LOCATION=prd \
  humpbacks/humpback-server
```

By default, Humpback Server will expose the UI over port `8100` and expose a API server over port `8101` for receiving agent report. 

Humpback Server has now been installed. you can log into your Humpback Server instance by opening a web browser and going to:

```
http://localhost:8100
```

## Deployment Hummpback Agent

By default, Humpback Agent will expose a API server over port `8018` for receiving Humpback Server call. 

```bash

docker run -d \
  --name=humpback-agent \
  --net=host \
  --restart=always \
  --privileged \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/lib/docker:/var/lib/docker \
  -e HUMPBACK_AGENT_API_BIND=0.0.0.0:8018 \
  -e HUMPBACK_SERVER_HOST={server-address}:8101 \
  -e HUMPBACK_VOLUMES_ROOT_DIRECTORY=/var/lib/docker \
  humpbacks/humpback-agent

```

Please replace `{server-address}` to the Humbpack Server IP.