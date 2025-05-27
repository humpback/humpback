```shell

docker run -d \
--name=humpback \
--net=host \
--restart=always \
-v /etc/localtime:/etc/localtime \
-v /var/lib/humpback:/workspace/data \
-e LOCATION=dev \
-e SITE_PORT=8300 \
-e BACKEND_PORT=8301 \
docker.io/humpbacks/humpback-server:develop

```