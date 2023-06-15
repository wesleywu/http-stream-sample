package main

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	http.HandleFunc("/stream", handleStreamVanilla)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	//s := g.Server()
	//s.BindHandler("/stream", handleStream)
	//s.SetPort(8080)
	//s.Run()
}

func handleStream(r *ghttp.Request) {
	ctx := r.Context()

	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")

	// 参考规范
	// https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "Client disconnected unexpectedly")
			return
		default:
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
			g.Log().Infof(ctx, "sent message %s", string(jsonBytes))
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

func handleStreamVanilla(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	// 参考规范
	// https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "Client disconnected unexpectedly")
			return
		default:
			testMessage := SampleMessage{
				Name:  "name" + gconv.String(i),
				Value: "value" + gconv.String(i),
			}
			// 不推荐使用 MarshalIndent，尽管不违反规范，规范中允许多行 Event，即 Event 消息体中允许单个换行符
			// 但某些 HTTP 客户端例如 PostMan 不支持多行 Event 的显示
			// jsonBytes, _ := json.MarshalIndent(testMessage, "", "  ")
			jsonBytes, _ := json.Marshal(testMessage)
			eventBuf := bytes.Buffer{}
			// 根据规范，Event 消息体的开头 需要有 "data: " 前缀
			eventBuf.WriteString("data: ")
			eventBuf.Write(jsonBytes)
			eventBuf.WriteByte('\n')
			eventBuf.WriteByte('\n')
			_, err := w.Write(eventBuf.Bytes())
			if err != nil {
				http.Error(w, "Failed to write to response", http.StatusInternalServerError)
				return
			}
			flusher.Flush()
			g.Log().Infof(ctx, "sent message %s", string(jsonBytes))
			time.Sleep(time.Millisecond * 1000)
		}
	}
}
