## Quick Start
### clone该仓库
```bash
git clone https://github.com/hts0000/icode-by-pass.git
```

### 解压并导入镜像
链接：https://pan.baidu.com/s/1vIhGSEySsP1Te2TkyXOLnQ  
提取码：bzax
```bash
# 解压镜像
tar -xf coolenv-icode-bypass.tar.gz
# 导入镜像
docker load < coolenv-icode-bypass
```

### 启动https server
```bash
# 进入仓库目录
cd icode-by-pass
# 启动http server
go run go-test/server/main.go
```

### 启停容器
```bash
# 首次启动
# --add-host 修改为自己的 docker 网关ip
docker run --entrypoint start.sh --name coolenv -p 18000:18000 --add-host apis.imooc.com:172.17.0.1 -e ICODE="xxxxxx" -e GODEBUG=x509ignoreCN=0 coolenv-icode-bypass

# 后续启动
docker start coolenv

# 停止容器
docker kill coolenv
```
**注意: 后续启动同样需要先启动`https server`**

## 原理
docker镜像本质是一堆文件的集合，既然是文件，那就可以直接访问。
```bash
# 查看镜像文件存放的位置
docker inspect -f "lowerdir={{.GraphDriver.Data.LowerDir}},upperdir={{.GraphDriver.Data.UpperDir}},workdir={{.GraphDriver.Data.WorkDir}}" <镜像名> 
```

我们可以把镜像挂载到本地，方便后续操作
```bash
mkdir -p /tmp/coolenv
mount -t overlay overlay -o $(docker inspect -f "lowerdir={{.GraphDriver.Data.LowerDir}},upperdir={{.GraphDriver.Data.UpperDir}},workdir={{.GraphDriver.Data.WorkDir}}" <镜像名>) /tmp/coolenv
```
分析一下镜像启动流程
```bash
# 查看镜像的启动流程
docker history --no-trunc <镜像名>
```
分析后确定是start.sh这个脚本负责项目的启动，执行`find /tmp/coolenv -name start.sh`命令，查找start.sh脚本。脚本内容如下：
```bash
#!/bin/bash
set -eu

echo Starting mongod...
entrypoint-mongo.sh mongod &
echo Starting rabbitmq...
entrypoint-rabbit.sh rabbitmq-server &
echo Starting api server...
coolenv --race_data=/data/coolenv/ningbo.json --grpc_pb_gen=/data/coolenv/pb/coolenv.pb.go --grpc_v2_pb_gen=/data/coolenv/pb/v2
```

可以发现认证是通过coolenv这个二进制文件来做的

继续执行命令`find /tmp/coolenv -name coolenv`查找coolenv，拿到二进制文件进行分析（具体分析过程Baidu），发现一个疑似的api网址：[https://apis.imooc.com/?cid=108&icode=%sinternal]()，浏览器请求发现返回一串json，把内容的乱码拿去解析发现内容为**icode不正确**。可以确定这个就是coolenv启动时请求的icode api。

既然知道了他请求的网址，那就可以在本地伪造一个fake imooc api网站，把镜像的请求截取到本地，直接返回相应内容，即可让coolenv认为认证成功。

首先要知道的是请求的这个网站是https的，而https的网站需要解决证书与认证的问题。

所以我们需要在本地生成一对根证书和服务端证书。生成证书步骤Baidu，在**仓库ssl目录**下有已经生成好的证书，可以直接使用。重点是生成时需要在生成配置中，配置Common Name必须为api.imooc.com
![](https://cdn.jsdelivr.net/gh/hts0000/images/2.jpg)

参考文章：[https://github.com/flyingtime/go-https]()

有了这一对根证书和服务端证书后，我们还需要把根证书写入镜像，并在本地启动https服务端，证书用上面生成的服务端证书。https服务端可以用go简单写一个，在go-test/server中可以找到。

根证书写入镜像比较麻烦，我的做法是在本机先写入ca文件，然后把ca文件拷贝到镜像对应目录中，具体可参考这篇文章：[https://www.cnblogs.com/jiaoyiping/p/6629442.html]()

简单来说就是在本地把根证书拷贝到/usr/local/share/ca-certificates目录，然后执行sudo update-ca-certificates，然后把自动生成的/etc/ssl/certs/ca-certificates.crt 文件拷贝到镜像/etc/ssl/certs/目录中。

```bash
# 把镜像umount
cd /tmp
umount /tmp/coolenv
```

做完这一切之后，镜像就变成了一个认可我们本地https服务端的镜像，现在需要做的就是把镜像对[https://api.imooc.com]()的请求改为请求本地即可。

要做到这一点很简单，docker命令本身提供了--add-host选项可以为容器添加域名解析。

到这里就完成了所有工作，启动本地的https服务端，然后启动镜像，启动命令如下：
```bash
# --add-host: 改为自己的docker的网关ip
docker run --name coolenv -p 18000:18000 --add-host apis.imooc.com:172.17.0.1 -e ICODE="xxxxxx" -e GODEBUG=x509ignoreCN=0 镜像名
```

后续每次启动都需要先启动本地的https服务端。
