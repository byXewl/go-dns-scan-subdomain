基于go语言的DNS探测子域名

>DNS查询：对字典文件中的每个条目进行DNS查询，看它们是否解析到有效的IP地址。
当你发起一个DNS查询时，DNS服务器会检查其数据库，看是否有对应的记录。如果有，它会返回相应的IP地址；如果没有，它会返回一个NXDOMAIN（域名不存在）响应。
通过自动化工具对大量潜在的子域名进行查询，可以发现那些未被广泛知晓但实际上已经配置了DNS记录的子域名和其IP。

>使用DNS探测而不是用HTTP请求：域名解析的IP服务不一定支持HTTP协议。

对应以下命令：
```
windows:
nslookup -type=A example.com
-type=A指定查询A记录，即域名对应的IPv4地址。

host example.com

linux:
dig example.com
```