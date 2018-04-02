# Zookeeper Cluster Deployment

> build `Zookeeper` cluster   

- Download Zookeeper  
  
Using the `Zookeeper 3.4.11` as example. Please make sure the `java open jdk7` runtime environment is installed on the server. 

```bash
$ cd /opt/app
$ wget http://mirror.bit.edu.cn/apache/zookeeper/zookeeper-3.4.11/zookeeper-3.4.11.tar.gz
$ tar -xzvf zookeeper-3.4.11.tar.gz
$ mv zookeeper-3.4.11 zookeeper
$ cd zookeeper
```
- Modify the configuration file   

```bash
$ cp conf/zoo_sample.cfg conf/zoo.cfg
$ vim conf/zoo.cfg
```

- The configuration file is modified as below 

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

- Create a data folder   

```bash
mkdir -p /opt/app/zookeeper/zkdata
mkdir -p /opt/app/zookeeper/logs
```

Do the same installation to the other servers which need install zookeeper. 

Next we will build `zookeeper` cluster.

- Create a Zookeeper node identification file `myid`   

Create `myid` number, execution on each server, pay attention to each server` myid` to correspond to the correct number.

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

- Start Zookeeper

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

- Zookeeper status check   

```bash
192.168.2.80
$ /opt/app/zookeeper/bin/zkServer.sh status
```
If the Zookeeper start failed, please check the below reasons:   

1、`Zoo.cfg` file configuration error: dataLogDir specified directory is not created.   

2、The integer in the `myid` file is not formatted correctly or does not correspond to the server integer in` zoo.cfg`.   

3、The firewall does not open the Zookeeper service port, such as `2888` and` 3888`.  
