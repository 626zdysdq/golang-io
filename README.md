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

![]()

我们可以将一个web服务器的工作流程总结为下面的过程：

* 客户机通过TCP/IP协议建立到服务器的TCP连接
* 客户端向服务器发送HTTP协议请求包，请求服务器里的资源文档
* 服务器向客户机发送HTTP协议应答包，如果请求的资源包含有动态语言的内容，那么服务器会调用动态语言的解释引擎负责处理“动态内容”，并将处理得到的数据返回给客户端
* 客户机与服务器断开。由客户端解释文档，在客户端屏幕上渲染显示。

关于web工作模式更详细的介绍请[点击这里](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.1.md)


#### 搭建一个web服务器

>* Go语言里面提供了一个完善的net/http包，通过http包可以很方便的搭建起来一个可以运行的Web服务。同时使用这个包能很简单地对Web的路由，静态文件，模版，cookie等数据进行设置和操作。

* 下面是实现的一个简单的web服务器

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

* 针对上面的简单实现，我们需要了解的两个地方
  * `http.HandleFunc`注册了请求的路由规则。如在原代码中，当请求uri为"/"，路由就会转到函数sayhelloName去执行。这在后面实现不同路径访问时是很重要的。
  * `http.ListenAndServe`处设置服务器监听端口，利用这个接口，我们就可以简单的运行起一个服务器
  

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


[源代码地址](https://github.com/626zdysdq/golang-io)