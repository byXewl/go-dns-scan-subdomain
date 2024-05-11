package main

import (
	"fmt"
	"github.com/sechelper/seclib/async"   // 引入异步处理库
	"github.com/sechelper/seclib/dict"    // 引入字典处理库
	"github.com/sechelper/seclib/network" // 引入网络库，用于DNS查询
	"log"                                 // 用于记录日志
	"os"                                  // 用于文件操作
	"time"                                // 用于处理时间相关的操作
)

// Lookdict 函数用于读取字典文件并打印每一行
func Lookdict() {
	path := "subdomain.txt"  // 字典文件的路径
	op, err := os.Open(path) // 打开字典文件
	if err != nil {
		log.Fatal(err) // 如果打开文件失败，记录错误并终止程序
	}
	defer op.Close() // 确保文件最终被关闭

	d := dict.NewDict(op) // 使用文件句柄创建字典对象
	if err != nil {
		log.Fatal(err) // 如果创建字典失败，记录错误并终止程序
	}

	for d.Scan() { // 遍历字典中的所有条目
		if line, err := d.Line(); err == nil {
			fmt.Println(line) // 打印当前条目
		}
	}
}

func main() {
	// 调用 Lookdict 函数查看字典文件内容
	// Lookdict()

	domain := "doubao.com" // 定义要查询的域名

	// 配置DNS解析器
	resolver := network.Dns{
		NewMsg:  network.NewDefaultMsg, // 使用默认的消息生成器
		Ns:      "8.8.8.8",             // 设置DNS服务器地址
		Timeout: 5 * time.Second,       // 设置查询超时时间
	}

	// 使用 async.Goroutine 启动指定数量的goroutine进行并发处理
	async.Goroutine(10, func(c *chan any) { // 第一个参数是goroutine的数量
		path := "subdomain.txt"  // 字典文件路径
		op, err := os.Open(path) // 打开字典文件
		if err != nil {
			log.Fatal(err) // 如果打开文件失败，记录错误并终止程序
		}
		defer op.Close() // 确保文件最终被关闭

		dt := dict.NewDict(op) // 创建字典对象
		if err != nil {
			log.Fatal(err) // 如果创建字典失败，记录错误并终止程序
		}

		for dt.Scan() { // 遍历字典中的所有条目
			*c <- dt.Text() // 将条目发送到通道
		}
	}, func(a ...any) { // 这里的 a 是从通道中读取的字典条目
		subdomain := a[0].(string)                              // 获取字典中的子域名
		ips, err := resolver.LookupIP(subdomain + "." + domain) // 对子域名进行DNS查询
		if err != nil {
			return // 如果查询失败，忽略当前条目
		}
		if len(ips) != 0 {
			fmt.Println("[*]", subdomain+"."+domain, ips) // 如果查询成功，打印结果
		} else {
			fmt.Println("[x]", subdomain+"."+domain) // 如果查询无结果，打印标记
		}
	})
}
