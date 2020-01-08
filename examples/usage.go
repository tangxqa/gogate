package main

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/tangxqa/gogate/conf"
	serv "github.com/tangxqa/gogate/server"
)

func main() {
	// 初始化
	serv.InitGogate("gogate.yml", "log.xml")

	// 创建Server
	server, err := serv.NewGatewayServer(
		conf.App.ServerConfig.Host,
		conf.App.ServerConfig.Port,
		conf.App.EurekaConfig.RouteFile,
		conf.App.ServerConfig.MaxConnection,
		// 是否启用优雅关闭
		true,
		// 优雅关闭最大等待时间, 上一个参数为true时有效
		time.Second * 30,
	)
	if nil != err {
		fmt.Println(err)
		return
	}

	// ******************* 非必须 *************************
	// 注册自定义过虑器, 在转发请求之前调用
	customPreFilter := serv.NewPreFilter("pre-log-filter1", PreLogFilter)
	server.AppendPreFilter(customPreFilter)
	// 在指定filter后面添加指定filter
	server.InsertPreFilterBehind("pre-log-filter1", customPreFilter)
	fmt.Printf("pre filters: %v\n", server.ExportAllPreFilters())

	// optional
	// 注册自定义过虑器, 在转发请求之后调用
	customPostFilter := serv.NewPostFilter("post-log-filter1", PostLogFilter)
	server.AppendPostFilter(customPostFilter)
	// 在指定filter后面添加指定filter
	server.InsertPostFilterBehind("pre-log-filter1", customPostFilter)
	fmt.Printf("post filters: %v\n", server.ExportAllPostFilters())

	// 自定义过虑器的添加方法必须在server启动之前调用, 启动后调用无效
	// ******************* 非必须 *************************

	// 启动Server
	err = server.Start()
	if nil != err {
		fmt.Println(err)
		return
	}

	// 等待优雅关闭
	server.WaitForGracefullyClose()
}

// 此方法会在gogate转发请求之前调用
// server: gogate服务器对象
// ctx: fasthttp请求上下文
// newRequest: gogate在转发请求时使用的请求对象指针, 可以做一些修改, 比如改请求参数,添加请求头之类
// return: 返回true则会继续执行后面的过虑器(如果有的话), 返回false则不会执行
func PreLogFilter(server *serv.Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	fmt.Println("request path: " + ctx.URI().String())

	return true
}

// 此方法会在gogate转发请求之后调用
// req: 转发给上游服务的HTTP请求
// resp: 上游服务的响应
// return: 返回true则会继续执行后面的过虑器(如果有的话), 返回false则不会执行
func PostLogFilter(req *fasthttp.Request, resp *fasthttp.Response) bool {
	fmt.Println("response: " + resp.String())

	return true
}
