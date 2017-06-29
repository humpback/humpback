# Deploy The Humpback Site

> Humpback Web    

We will use `192.168.2.80:8012` as the Humpback Web site.

Site default port: 80, if you want to customize the port, please specify `-e HUMPBACK_LISTEN_PORT = XXXX` when the container is created.   

`/opt/app/humpback-web/dbFiles` is the Humpback system persistent data files, will store site management and grouping information.

```bash
$ ssh root@192.168.2.80
$ mkdir -p /opt/app/humpback-web
$ docker run -d --net=host --restart=always \
 -e HUMPBACK_LISTEN_PORT=8012 \
 -v /opt/app/humpback-web/dbFiles:/humpback-web/dbFiles \
 --name humpback-web \
 humpbacks/humpback-web:1.0.0
```

Open browser: http://192.168.2.80:8012    

Default Account: admin Password: 123456   

![Humpback Web](_media/humpback-web.png)

Log in to Humpback as an administrator and expand the left `Manage` and click `System Config`. Check the `Enable Private Registry`, and fill in the previously deployed private registry service address `http: //192.168.2.80: 5000` and `Save`. 

![Humpback Web Registry Bind](_media/humpback-web-registry-bind.png)
