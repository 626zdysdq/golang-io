---
title: GO语言处理web程序的输入输出
data: 2018-11-14 23:55:59
tags:
    - 服务计算
    - Go语言
    - web service
    
---

## 概述
设计一个 web 小应用，展示静态文件服务、js 请求支持、模板输出、表单处理、Filter 中间件设计等方面的能力。（不需要数据库支持）

## 任务要求
1. 支持静态文件服务
2. 支持简单 js 访问
3. 提交表单，并输出一个表格
4. 对 `/unknown` 给出开发中的提示，返回码 `5xx`

<!--more-->

## 服务器的工作原理

>* 一个Web服务器也被称为HTTP服务器，它通过HTTP协议与客户端通信。在实现一个服务器之前，我们首先要了解一下服务器的工作原理

参考老师博客Web服务的工作模式流程图

![模式流程](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8810.28.32.png?raw=true)

我们可以将一个web服务器的工作流程总结为下面的过程：

* 客户机通过TCP/IP协议建立到服务器的TCP连接
* 客户端向服务器发送HTTP协议请求包，请求服务器里的资源文档
* 服务器向客户机发送HTTP协议应答包，如果请求的资源包含有动态语言的内容，那么服务器会调用动态语言的解释引擎负责处理“动态内容”，并将处理得到的数据返回给客户端
* 客户机与服务器断开。由客户端解释文档，在客户端屏幕上渲染显示。

关于web工作模式更详细的介绍请[点击这里](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.1.md)

## 任务实现

### 搭建一个web服务器

>* Go语言里面提供了一个完善的net/http包，通过http包可以很方便的搭建起来一个可以运行的Web服务。同时使用这个包能很简单地对Web的路由，静态文件，模版，cookie等数据进行设置和操作。

* 下面是实现的一个简单的web服务器。运行下面的代码，在浏览器中输入`http://localhost:9090/`就可以在网页上显示`HelloWorld!`

```
package main
import (
    "fmt"
    "net/http"
    "strings"
    "log"
)
func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm() //解析参数，默认是不会解析的
    fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello Wrold!") //这个写入到w的是输出到客户端的
}
func main() {
    http.HandleFunc("/", sayhelloName) //设置访问的路由
    err := http.ListenAndServe(":9090", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

```

* 针对上面的简单实现，我们会发现实现一个GO语言的服务器只需要几行代码，其中最关键的是包括一下两个部分

#### 注册路由规则

>* `http.HandleFunc`注册了请求的路由规则。如在原代码中，当请求uri为"/"，路由就会转到函数sayhelloName去执行。这在后面实现不同路径访问时是很重要的。其源码实现如下：

```
// HandleFunc registers the handler function for the given pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}

```

根据源码我们可以看到，`HandleFunc`将传入的处理响应函数与对应的path进行匹配，之后调用`mux.Handle`对传入的路径进行解析，然后向ServeMux中添加路由规则。而其完整的参数传递过程为`mux.ServerHTTP->mux.Handler->mux.handler->mux.match`。对于这部分实现有兴趣的同学，可以去阅读实现过程的源码。如果只需要使用，只要记住`http.HandleFunc`就可以的参数传递，使用方式即可。

>* 在本任务的实现中，我利用了改函数创建了静态文件、未知路径访问以及表单提交的不同路由规则，其源代码和实现过程如下：

```
//对静态文件，js文件的请求路由规则
mx.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/add/"))))

//对表单的请求路由规则
mx.HandleFunc("/", home).Methods("GET")
mx.HandleFunc("/", login).Methods("POST")

//对错误访问的请求路由规则
mx.HandleFunc("/api/test", apiTestHandler(formatter)).Methods("GET")
mx.HandleFunc("/unknown", NotImplemented).Methods("GET")

```

其实现效果如下：

![路由规则](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.03.27.png)



#### 设置监听函数

>* `http.ListenAndServe`处设置服务器监听端口，利用这个接口，我们就可以简单的运行起一个服务器，开始监听请求信息，其实现的源代码如下：

```
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
  
在源代码中我们可以看到，该函数首先实例化了`Server`，接着调用了`Server.ListenAndServe()`，其源码如下：

```
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)    //监听端口
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

在其中，最重要的就是`net.Listen("tcp",addr)`部分，和`srv.Serve()`部分。这两个函数分别设置了该服务器的监听端口以及处理提交向服务器的请求。特别是`srv.Serve()`，它的实现过程确保了每个请求都能保持独立，相互不会阻塞，以达到可以高效响应网络事件的目的
  
  

## 实现效果

#### 访问静态文件


* 访问静态的图片

![金木](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.03.53.png?raw=true)

* 访问静态图标

![圆](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.03.50.png?raw=true)

* 访问静态网页(html+css+js)

![html1](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.03.59.png?raw=true)

![html2](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.04.04.png?raw=true)

#### 访问js文件

![js](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.03.41.png?raw=true)


#### 提交表单，并输出表格
>* 提交表单

![up](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.42.34.png?raw=true)

>* 输出表格

![show](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%881.29.14.png?raw=true)

#### 对未知请求的错误处理

![wrong](https://github.com/626zdysdq/golang-io/blob/master/picture/%E5%B1%8F%E5%B9%95%E5%BF%AB%E7%85%A7%202018-11-16%20%E4%B8%8A%E5%8D%8812.49.40.png?raw=true)

#### 服务器的终端显示

![console]()

[源代码地址](https://github.com/626zdysdq/golang-io)