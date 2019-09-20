# 仓库

### 富文本编辑器

TinyMCE: [https://www.tiny.cloud/docs/demo/full-featured/#](https://www.tiny.cloud/docs/demo/full-featured/#)

CKEditor5: [https://ckeditor.com/ckeditor-5/demo/](https://ckeditor.com/ckeditor-5/demo/)

--------

### KCP

```bash
./client_linux_amd64 -r "<remote_ip>:<server_port>" -l ":<client_port>" -mode fast3 -nocomp -autoexpire 900 -sockbuf 16777217 -dscp 46
./server_linux_amd64 -t "localhost:21889" -l ":<server_port>" -mode fast3 -nocomp -sockbuf 16777217 -dscp 46
```

--------

### iptables 端口转发

```bash
iptables -t nat -A PREROUTING -p udp --dport <local_port> -j DNAT --to-destination <remote_ip>:<remote_port>
iptables -t nat -A POSTROUTING -p udp -d <remote_ip> --dport <remote_port> -j SNAT --to-source <内网IP>
```

--------

国外黑名单：[SS pac.txt](https://drive.google.com/open?id=1iXrz5LLodiO85tUqKvtwhGrFSYqxteeq)

跨浏览器测试：[https://www.browserling.com/](https://www.browserling.com/) 免费套餐每次限时3分钟，月付套餐$19。

傅里叶级数画图：[https://www.youtube.com/watch?v=ds0cmAV-Yek](https://www.youtube.com/watch?v=ds0cmAV-Yek)

Redis 文档参考: [http://redisdoc.com/index.html](http://redisdoc.com/index.html)

Spring AOP ProceedingJoinPoint 获取执行函数类名和方法名：
```java
String className = joinPoint.getTarget().getClass().getName();
String methodName = joinPoint.getSignature().getName();
```

--------

### Turn 服务安装

```bash
wget -O turn.tar.gz http://turnserver.open-sys.org/downloads/v4.5.0.3/turnserver-4.5.0.3.tar.gz 

# Download the source tar

tar -zxvf turn.tar.gz

 # unzip

cd turnserver-*

./configure

make && make install
```

```bash
export USERS='test=test'
export REALM=my-server.com
export UDP_PORT=3478
```

运行服务

`sudo turnserver -a -o -v -n  --no-dtls --no-tls -u test:test -r "someRealm"`

--------

### 兼容获取 `getUserMedia`
`window.navigator.getUserMedia = navigator.getUserMedia || navigator.webKitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia;`

--------

### 命令行窗口代理

Windows:
```bash
set http_proxy=http://127.0.0.1:1080
set https_proxy=http://127.0.0.1:1080
```

MacOS/Linux:
```bash
export http_proxy=http://127.0.0.1:1087
export https_proxy=http://127.0.0.1:1087
```

--------

```bash
wget --no-check-certificate -O shadowsocks-all.sh https://raw.githubusercontent.com/AppSo/shadowsocks_install/master/shadowsocks-all.sh
chmod +x shadowsocks-all.sh
./shadowsocks-all.sh 2>&1 | tee shadowsocks-all.log
```

--------

### 线路测评

厂商 | 线路 | 评价 | 测试
---|---|---|---
HostKVM |香港CN2优化 | 延迟低，带宽太小，访问速度慢 | - 
GbpsCloud | 无锡日本专线 | 延迟较低40ms，带宽适中10Mbps，价格高￥120/月 | https://www.speedtest.net/zh-Hans/result/8486106593
CloudIPLC | 上海日本专线 | 延迟极低17ms，带宽较好20Mbps | https://www.speedtest.net/zh-Hans/result/8486104142
BWG | 美国CN2 GIA | 延迟一般 ~170ms 带宽充裕>50Mbps | https://www.speedtest.net/result/8486108890

--------

### GOPROXY公共服务

1. https://mirrors.aliyun.com/goproxy/
2. https://goproxy.cn

设置GOPROXY环境变量
`export GOPROXY=https://goproxy.foobar.com/`

--------

### Gson反序列化List

```java
Type type = new TypeToken<List<FooBarClass>>(){}.getType();
List<FooBarClass> list = new Gson().fromJson(jsonString, type);
```

--------

### 取消Chrome浏览器的“托管状态”

- 打开注册表编辑器（右击左下角Windows徽标按钮-运行-输入 `regedit`-确定 ）
- 定位至 `计算机\HKEY_LOCAL_MACHINE\SOFTWARE\Policies\Google\Chrome`
- 删除左侧Chrome下的 `EnabledPlugins`
- 重启Chrome

--------

### MySQL 连接命令

`mysql -h主机地址 -u用户名 -p用户密码`

例如连接本机 root `mysql -uroot -p` 输入密码登录。

代码抽象一致性、圈复杂度低、以逻辑变量代替大段逻辑语句、以多个逻辑变量的卫语句实现分解的逻辑链路，五个1原则。

--------

```sql
select * from test where id = 100
```

```xml
<select id="getEmployeesListParams" resultType="Employees">
        select *
        from EMPLOYEES e
        where e.EMPLOYEE_ID in
        <foreach collection="list" item="employeeId" index="index"
            open="(" close=")" separator=",">
            #{employeeId}
        </foreach>
</select>
```

> hello