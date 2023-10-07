# tRPC Consul 名字服务

## 关于 consul

见 [HashiCorp - Consul](https://www.consul.io/) 

## 示例
配置：
```yaml
plugins:
  naming:
    consul:
      address: dev.cloud.com:8500
      services:
        - trpc.test.helloworld.Greeter  # 一定要与 trpc service 相同
      register:  #  默认注册配置，上面的 services 会使用
        interval: 1s
        timeout: 1s
        tags:
          - test
        meta:
          appid: 1
        weight: 10
        deregister_critical_service_after: 10m
      services_register:  # 独立注册配置，不同服务可以有不同配置
        - service: trpc.test.helloworld.Greeter  # 一定要与 trpc service 相同
          register:  #  默认注册配置，上面的 services 会使用
            interval: 1s
            timeout: 1s
            tags:
              - test
            meta:
              appid: 1
            weight: 10
            deregister_critical_service_after: 10m
      selector:
        loadBalancer: random

client:  # 客户端调用的后端配置
  service:  # 针对单个后端的配置
    - callee: trpc.test.helloworld.Greeter  # 后端服务协议文件的 service name, 如何 callee 和下面的 name 一样，那只需要配置一个即可
      target: consul://trpc.test.helloworld.Greeter  # 后端服务地址 consul
      network: tcp  # 后端服务的网络类型 tcp udp
      protocol: http  # 应用层协议 trpc http
      timeout: 10000  # 请求最长处理时间
      serialization: 2  # 序列化方式 0-pb 1-jce 2-json 3-flatbuffer，默认不要配置

```

main 入口：
```go
package main

import (
    "trpc.group/trpc-go/trpc-go"
    "trpc.group/trpc-go/trpc-go/server"
    _ "trpc.group/trpc-go/trpc-naming-consul"
)

func main() {
    s := trpc.NewServer()
    // do sth ...
    
    s.Serve()
}
```
