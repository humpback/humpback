# Zookeeper 集群部署

> `Zookeeper` 集群搭建   

- Zookeeper 下载   
  
  示例采用 `Zookeeper 3.4.6` 版本，请确保在服务器上安装了 `java open jdk7` 运行环境，下载 `Zookeeper` 并解压。

```bash
$ cd /opt/app
$ wget http://mirror.bit.edu.cn/apache/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz
$ tar -xzvf zookeeper-3.4.6.tar.gz
$ mv zookeeper-3.4.6 zookeeper
$ cd zookeeper
```
- 修改配置文件   

```bash
$ cp conf/zoo_sample.cfg conf/zoo.cfg
$ vim conf/zoo.cfg
```

- 配置文件修改后如下   

```
tickTime=2000
initLimit=10
syncLimit=5
dataDir=/opt/app/zookeeper/zkdata
dataLogDir=/opt/app/zookeeper/logs
clientPort=2181
server.1=192.168.2.80:2888:3888
server.2=192.168.2.81:2888:3888
server.3=192.168.2.82:2888:3888
```

- 创建数据目录   

```bash
mkdir -p /opt/app/zookeeper/zkdata
mkdir -p /opt/app/zookeeper/logs
```

&ensp;&ensp;&ensp;以上安装流程三台服务器都照执行，接下来开始创建 `zookeeper` 集群并启动。   

- 建立 Zookeeper 节点标识文件 `myid`   

&ensp;&ensp;&ensp;创建 `myid` 编号，依次在每台 Server 上执行，注意每台 Server 的 `myid` 要对应正确的编号。

```bash
192.168.2.80
$ echo "1" > /opt/app/zookeeper/zkdata/myid
```

```bash
192.168.2.81  
$ echo "2" > /opt/app/zookeeper/zkdata/myid
```

```bash
192.168.2.82
$ echo "3" > /opt/app/zookeeper/zkdata/myid
```

- 启动 Zookeeper

```bash
192.168.2.80
$ /opt/app/zookeeper/bin/zkServer.sh start
```

```bash
192.168.2.81  
$ /opt/app/zookeeper/bin/zkServer.sh start
```

```bash
192.168.2.82
$ /opt/app/zookeeper/bin/zkServer.sh start
```

- Zookeeper 状态检查   

```bash
192.168.2.80
$ /opt/app/zookeeper/bin/zkServer.sh status
```
&ensp;&ensp;&ensp;如果 Zookeeper 启动失败或报错，可能有以下几个原因：   

&ensp;&ensp;&ensp;2、`zoo.cfg` 文件配置出错：dataLogDir 指定的目录未被创建。   

&ensp;&ensp;&ensp;3、`myid` 文件中的整数格式不对，或者与 `zoo.cfg` 中的 server 整数不对应。   

&ensp;&ensp;&ensp;4、防火墙未打开 Zookeeper 服务端口，如 `2888` 和 `3888`。  