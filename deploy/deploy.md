```shell

## 使用本地目录挂载

docker run -d \
--name=humpback \
--net=host \
--restart=always \
-v /etc/localtime:/etc/localtime \
-v /var/lib/humpback/data:/workspace/data \
-v /var/lib/humpback/certs:/workspace/certs \
-e SITE_PORT=8300 \
-e BACKEND_PORT=8301 \
docker.io/humpbacks/humpback:latest


## 使用volune挂载

docker run -d \
--name=humpback \
--net=host \
--restart=always \
-v /etc/localtime:/etc/localtime \
-v humpback_data:/workspace/data \
-v humpback_certs:/workspace/certs \
-e SITE_PORT=8300 \
-e BACKEND_PORT=8301 \
docker.io/humpbacks/humpback:latest

```