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

		for i := 0; i < 10; i++ {
			testMessage := SampleMessage{
				Name:  "name" + gconv.String(i),
				Value: "value" + gconv.String(i),
			}
			jsonBytes, _ := json.MarshalIndent(testMessage, "", "  ")
			r.Response.Writefln("data: " + string(jsonBytes) + "\n")
			r.Response.Flush()
			time.Sleep(time.Millisecond * 1000)
		}
	})
	s.SetPort(8080)
	s.Run()
}
