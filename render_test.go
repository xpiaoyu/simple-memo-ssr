package main

import (
	"io/ioutil"
	"log"
	"testing"
)

var html = `<h1>仓库</h1>

<h3>富文本编辑器</h3>

<p>TinyMCE: <a href="https://www.tiny.cloud/docs/demo/full-featured/#">https://www.tiny.cloud/docs/demo/full-featured/#</a></p>

<p>CKEditor5: <a href="https://ckeditor.com/ckeditor-5/demo/">https://ckeditor.com/ckeditor-5/demo/</a></p>

<hr />

<h3>KCP</h3>
<pre style="color:#f8f8f2;background-color:#282a36">./client_linux_amd64 -r <span style="color:#f1fa8c">&#34;&lt;remote_ip&gt;:&lt;server_port&gt;&#34;</span> -l <span style="color:#f1fa8c">&#34;:&lt;client_port&gt;&#34;</span> -mode fast3 -nocomp -autoexpire <span style="color:#bd93f9">900</span> -sockbuf <span style="color:#bd93f9">16777217</span> -dscp <span style="color:#bd93f9">46</span>
./server_linux_amd64 -t <span style="color:#f1fa8c">&#34;localhost:21889&#34;</span> -l <span style="color:#f1fa8c">&#34;:&lt;server_port&gt;&#34;</span> -mode fast3 -nocomp -sockbuf <span style="color:#bd93f9">16777217</span> -dscp <span style="color:#bd93f9">46</span>
</pre>
<hr />

<h3>iptables 端口转发</h3>
<pre style="color:#f8f8f2;background-color:#282a36">iptables -t nat -A PREROUTING -p udp --dport &lt;local_port&gt; -j DNAT --to-destination &lt;remote_ip&gt;:&lt;remote_port&gt;
iptables -t nat -A POSTROUTING -p udp -d &lt;remote_ip&gt; --dport &lt;remote_port&gt; -j SNAT --to-source &lt;内网IP&gt;
</pre>
<hr />

<p>国外黑名单：<a href="https://drive.google.com/open?id=1iXrz5LLodiO85tUqKvtwhGrFSYqxteeq">SS pac.txt</a></p>

<p>跨浏览器测试：<a href="https://www.browserling.com/">https://www.browserling.com/</a> 免费套餐每次限时3分钟，月付套餐$19。</p>

<p>傅里叶级数画图：<a href="https://www.youtube.com/watch?v=ds0cmAV-Yek">https://www.youtube.com/watch?v=ds0cmAV-Yek</a></p>

<p>Redis 文档参考: <a href="http://redisdoc.com/index.html">http://redisdoc.com/index.html</a></p>

<p>Spring AOP ProceedingJoinPoint 获取执行函数类名和方法名：</p>
<pre style="color:#f8f8f2;background-color:#282a36">String className <span style="color:#ff79c6">=</span> joinPoint<span style="color:#ff79c6">.</span><span style="color:#50fa7b">getTarget</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">.</span><span style="color:#50fa7b">getClass</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">.</span><span style="color:#50fa7b">getName</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">;</span>
String methodName <span style="color:#ff79c6">=</span> joinPoint<span style="color:#ff79c6">.</span><span style="color:#50fa7b">getSignature</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">.</span><span style="color:#50fa7b">getName</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">;</span>
</pre>
<hr />

<h3>Turn 服务安装</h3>
<pre style="color:#f8f8f2;background-color:#282a36">wget -O turn.tar.gz http://turnserver.open-sys.org/downloads/v4.5.0.3/turnserver-4.5.0.3.tar.gz 

<span style="color:#6272a4"># Download the source tar
</span><span style="color:#6272a4"></span>
tar -zxvf turn.tar.gz

 <span style="color:#6272a4"># unzip
</span><span style="color:#6272a4"></span>
<span style="color:#8be9fd;font-style:italic">cd</span> turnserver-*

./configure

make <span style="color:#ff79c6">&amp;&amp;</span> make install
</pre><pre style="color:#f8f8f2;background-color:#282a36"><span style="color:#8be9fd;font-style:italic">export</span> <span style="color:#8be9fd;font-style:italic">USERS</span><span style="color:#ff79c6">=</span><span style="color:#f1fa8c">&#39;test=test&#39;</span>
<span style="color:#8be9fd;font-style:italic">export</span> <span style="color:#8be9fd;font-style:italic">REALM</span><span style="color:#ff79c6">=</span>my-server.com
<span style="color:#8be9fd;font-style:italic">export</span> <span style="color:#8be9fd;font-style:italic">UDP_PORT</span><span style="color:#ff79c6">=</span><span style="color:#bd93f9">3478</span>
</pre>
<p>运行服务</p>

<p><code>sudo turnserver -a -o -v -n  --no-dtls --no-tls -u test:test -r &quot;someRealm&quot;</code></p>

<hr />

<h3>兼容获取 <code>getUserMedia</code></h3>

<p><code>window.navigator.getUserMedia = navigator.getUserMedia || navigator.webKitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia;</code></p>

<hr />

<h3>命令行窗口代理</h3>

<p>Windows:</p>
<pre style="color:#f8f8f2;background-color:#282a36"><span style="color:#8be9fd;font-style:italic">set</span> <span style="color:#8be9fd;font-style:italic">http_proxy</span><span style="color:#ff79c6">=</span>http://127.0.0.1:1080
<span style="color:#8be9fd;font-style:italic">set</span> <span style="color:#8be9fd;font-style:italic">https_proxy</span><span style="color:#ff79c6">=</span>http://127.0.0.1:1080
</pre>
<p>MacOS/Linux:</p>
<pre style="color:#f8f8f2;background-color:#282a36"><span style="color:#8be9fd;font-style:italic">export</span> <span style="color:#8be9fd;font-style:italic">http_proxy</span><span style="color:#ff79c6">=</span>http://127.0.0.1:1087
<span style="color:#8be9fd;font-style:italic">export</span> <span style="color:#8be9fd;font-style:italic">https_proxy</span><span style="color:#ff79c6">=</span>http://127.0.0.1:1087
</pre>
<hr />
<pre style="color:#f8f8f2;background-color:#282a36">wget --no-check-certificate -O shadowsocks-all.sh https://raw.githubusercontent.com/AppSo/shadowsocks_install/master/shadowsocks-all.sh
chmod +x shadowsocks-all.sh
./shadowsocks-all.sh <span style="color:#bd93f9">2</span>&gt;&amp;<span style="color:#bd93f9">1</span> | tee shadowsocks-all.log
</pre>
<hr />

<h3>线路测评</h3>

<table>
<thead>
<tr>
<th>厂商</th>
<th>线路</th>
<th>评价</th>
<th>测试</th>
</tr>
</thead>

<tbody>
<tr>
<td>HostKVM</td>
<td>香港CN2优化</td>
<td>延迟低，带宽太小，访问速度慢</td>
<td>-</td>
</tr>

<tr>
<td>GbpsCloud</td>
<td>无锡日本专线</td>
<td>延迟较低40ms，带宽适中10Mbps，价格高￥120/月</td>
<td><a href="https://www.speedtest.net/zh-Hans/result/8486106593">https://www.speedtest.net/zh-Hans/result/8486106593</a></td>
</tr>

<tr>
<td>CloudIPLC</td>
<td>上海日本专线</td>
<td>延迟极低17ms，带宽较好20Mbps</td>
<td><a href="https://www.speedtest.net/zh-Hans/result/8486104142">https://www.speedtest.net/zh-Hans/result/8486104142</a></td>
</tr>

<tr>
<td>BWG</td>
<td>美国CN2 GIA</td>
<td>延迟一般 ~170ms 带宽充裕&gt;50Mbps</td>
<td><a href="https://www.speedtest.net/result/8486108890">https://www.speedtest.net/result/8486108890</a></td>
</tr>
</tbody>
</table>

<hr />

<h3>GOPROXY公共服务</h3>

<ol>
<li><a href="https://mirrors.aliyun.com/goproxy/">https://mirrors.aliyun.com/goproxy/</a></li>
<li><a href="https://goproxy.cn">https://goproxy.cn</a></li>
</ol>

<p>设置GOPROXY环境变量
<code>export GOPROXY=https://goproxy.foobar.com/</code></p>

<hr />

<h3>Gson反序列化List</h3>
<pre style="color:#f8f8f2;background-color:#282a36">Type type <span style="color:#ff79c6">=</span> <span style="color:#ff79c6">new</span> TypeToken<span style="color:#ff79c6">&lt;</span>List<span style="color:#ff79c6">&lt;</span>FooBarClass<span style="color:#ff79c6">&gt;</span><span style="color:#ff79c6">&gt;</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">{</span><span style="color:#ff79c6">}</span><span style="color:#ff79c6">.</span><span style="color:#50fa7b">getType</span><span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">;</span>
List<span style="color:#ff79c6">&lt;</span>FooBarClass<span style="color:#ff79c6">&gt;</span> list <span style="color:#ff79c6">=</span> <span style="color:#ff79c6">new</span> Gson<span style="color:#ff79c6">(</span><span style="color:#ff79c6">)</span><span style="color:#ff79c6">.</span><span style="color:#50fa7b">fromJson</span><span style="color:#ff79c6">(</span>jsonString<span style="color:#ff79c6">,</span> type<span style="color:#ff79c6">)</span><span style="color:#ff79c6">;</span>
</pre>
<hr />

<h3>取消Chrome浏览器的“托管状态”</h3>

<ul>
<li>打开注册表编辑器（右击左下角Windows徽标按钮-运行-输入 <code>regedit</code>-确定 ）</li>
<li>定位至 <code>计算机\HKEY_LOCAL_MACHINE\SOFTWARE\Policies\Google\Chrome</code></li>
<li>删除左侧Chrome下的 <code>EnabledPlugins</code></li>
<li>重启Chrome</li>
</ul>

<hr />

<h3>MySQL 连接命令</h3>

<p><code>mysql -h主机地址 -u用户名 -p用户密码</code></p>

<p>例如连接本机 root <code>mysql -uroot -p</code> 输入密码登录。</p>

<p>代码抽象一致性、圈复杂度低、以逻辑变量代替大段逻辑语句、以多个逻辑变量的卫语句实现分解的逻辑链路，五个1原则。</p>

<hr />
<script lang="xxx">test throne transform task</script>
`

var md = `# 测试

参数 | 示例值 | 参数描述
--- | --- | ---
productId | IF | 产品ID 如: IF, IC, IH, TS, TF, T
date| 2019-7-5 | 交易日期

    t.Timestamp = fi.ModTime().UnixNano() / 1e6
    sort.Sort(ArticleList)
    if _, err := c.WriteString("success"); err != nil {
        log.Println("[ERROR]", err)
    }
`

var htmlBytes = []byte(html)
var mdBytes = []byte(md)

func BenchmarkHighlightKeyword(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HighlightKeyword(html, "t")
	}
	b.StopTimer()
}

func BenchmarkHighlightKeywordBytes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HighlightKeywordBytes(htmlBytes, []byte("tt"))
	}
	b.StopTimer()
}

func BenchmarkMarkdownToHtml(b *testing.B) {
	bytes, err := ioutil.ReadFile("article/warehouse.md")
	if err != nil {
		log.Println("[ERROR]", "读取失败")
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarkdownToHtml(bytes)
	}
	b.StopTimer()
}

func TestHighlightKeyword(t *testing.T) {
	t.Log(HighlightKeyword(html, "t"))
}

func TestHighlightKeywordBytes(t *testing.T) {
	t.Log(string(HighlightKeywordBytes(htmlBytes, []byte("t"))))
}

func TestHighlightKeyword2(t *testing.T) {
	t.Log(HighlightKeyword("ababbbbbbba<script>bbb</script>", "a"))
}

func TestHighlightKeywordBytes2(t *testing.T) {
	t.Log(HighlightKeywordBytes([]byte("ababbbbbbba<script>bbb</script>"), []byte("a")))
}

func TestRune(t *testing.T) {
	str := "你好"
	r := []rune(str)
	b := []byte(str)
	for _, v := range r {
		t.Logf("%X ", v)
	}
	t.Log()
	for _, v := range b {
		t.Logf("%X ", v)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(html)
	}
	b.StopTimer()
}
