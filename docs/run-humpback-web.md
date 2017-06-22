# Deploy The Humpback Site

> Humpback Web    

&ensp;&ensp;&ensp;With 192.168.2.80 as the Humpback Web site server to start, set the port: 8012.

&ensp;&ensp;&ensp;Site default access port: 80, if you want to customize the port, please specify `-e HUMPBACK_LISTEN_PORT = XXXX` when the container is created.   

&ensp;&ensp;&ensp;`/opt/app/humpback-web/dbFiles` file for the Humpback system persistent data files, will store site management and grouping information, please save after the start.

```bash
$ ssh root@192.168.2.80
$ mkdir -p /opt/app/humpback-web
$ docker run -d --net=host --restart=always \
 -e HUMPBACK_LISTEN_PORT=8012 \
 -v /opt/app/humpback-web/dbFiles:/humpback-web/dbFiles \
 --name humpback-web \
 humpbacks/humpback-web:1.0.0
```

&ensp;&ensp;&ensp;Access the site, open the browser input: http://192.168.2.80:8012    

&ensp;&ensp;&ensp;Default Account: admin Password: 123456   

![Humpback Web](_media/humpback-web.png)

&ensp;&ensp;&ensp;以管理员身份登录 Humpback Web 站点，展开左侧 `Manage` 功能点击 `System Config`，进入系统配置界面，将 `Enabel Private Registry` 勾选，并填入先前部署的私有仓库服务地址 `http://192.168.2.80:5000` Save 即可。 

![Humpback Web Registry Bind](_media/humpback-web-registry-bind.png)
