package main

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

type SampleMessage struct {
	Name  string
	Value string
}

func main() {
	s := g.Server()
	s.BindHandler("/stream", func(r *ghttp.Request) {
		r.Response.Header().Set("Content-Type", "text/event-stream")
		r.Response.Header().Set("Cache-Control", "no-cache")
		r.Response.Header().Set("Connection", "keep-alive")

		// 参考规范
		// https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events
		for i := 0; i < 10; i++ {
			testMessage := SampleMessage{
				Name:  "name" + gconv.String(i),
				Value: "value" + gconv.String(i),
			}
			// 不推荐使用 MarshalIndent，尽管不违反规范，规范中允许多行 Event，即 Event 消息体中允许单个换行符
			// 但某些 HTTP 客户端例如 PostMan 不支持多行 Event 的显示
			// jsonBytes, _ := json.MarshalIndent(testMessage, "", "  ")
			jsonBytes, _ := json.Marshal(testMessage)
			// 根据规范，Event 消息体的开头 需要有 "data: " 前缀
			r.Response.Writef("data: " + string(jsonBytes) + "\n\n")
			r.Response.Flush()
			time.Sleep(time.Millisecond * 1000)
		}
	})
	s.SetPort(8080)
	s.Run()
}
