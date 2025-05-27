# Install Humpback

## Prerequisites

Humpback consists of two components, the Humpback, and the Humpback Agent. Both them run as lightweight Docker containers on a Docker engine.

To get started, you will need the latest version of Docker installed and working. We recommend following the [official installation instructions](https://docs.docker.com/engine/install/) for Docker - in particular, we advise against installing Docker via snap on Ubuntu distributions as you may run into compatibility issues.

## Deployment Humpback

First, create the volume that Humpback will use to store its database:

```bash
docker volume create humpback_data
```

Then, install the Humpback container:

```bash
docker run -d \
  --name humpback \
  -p 8100:8100 \
  -p 8101:8101 \
  --restart=always \
  -v humpback_data:/workspace/data \
  -e LOCATION=prd \
  humpbacks/humpback:latest
```

By default, Humpback will expose the UI over port `8100` and expose a API server over port `8101` for receiving agent report. 

Humpback has now been installed. you can log into your Humpback instance by opening a web browser and going to:

```
http://localhost:8100
```

You can use the account initialized by the system to log in. Both the username and the password are `humpback`. 

## Deployment Humpback Agent

By default, Humpback Agent will expose a API server over port `8018` for receiving Humpback call. 

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

## Add Node in Humpback

Open Humpback website, under the Administration menu, on the Nodes page, you can add the node that you just install agent on by IP address.