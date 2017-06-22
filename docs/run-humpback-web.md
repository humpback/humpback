# 部署 Humpback 站点

> Humpback Web    

&ensp;&ensp;&ensp;以 192.168.2.80 作为 Humpback Web 站点服务器启动，设置端口：8012。

&ensp;&ensp;&ensp;站点默认访问端口为：80，如果要自定义端口，请在容器创建时指定 `-e HUMPBACK_LISTEN_PORT=XXXX` 即可。   

&ensp;&ensp;&ensp;`/opt/app/humpback-web/dbFiles` 文件为 Humpback 系统持久化数据文件，会存储站点管理和分组信息，启动后请妥善保存。

```bash
$ ssh root@192.168.2.80
$ mkdir -p /opt/app/humpback-web
$ docker run -d --net=host --restart=always \
 -e HUMPBACK_LISTEN_PORT=8012 \
 -v /opt/app/humpback-web/dbFiles:/humpback-web/dbFiles \
 --name humpback-web \
 humpbacks/humpback-web:1.0.0
```

&ensp;&ensp;&ensp;访问站点，打开浏览器输入：http://192.168.2.80:8012    

&ensp;&ensp;&ensp;默认账户：admin  密码：123456   

![Humpback Web](_media/humpback-web.png)

&ensp;&ensp;&ensp;以管理员身份登录 Humpback Web 站点，展开左侧 `Manage` 功能点击 `System Config`，进入系统配置界面，将 `Enabel Private Registry` 勾选，并填入先前部署的私有仓库服务地址 `http://192.168.2.80:5000` Save 即可。 

![Humpback Web Registry Bind](_media/humpback-web-registry-bind.png)