# HTTP Streaming example in Golang 

## server 
在 HTTP 服务端发送 Stream response，每一个 Event 为一个 json 结构体

## client
在客户端通过 POST(用GET亦可) 访问服务端，对结果进行解码

## 参考规范
https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events