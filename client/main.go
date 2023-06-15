package main

import (
	"context"
	"io"
	"net/http"

	"github.com/WesleyWu/http-stream-sample/client/streamio"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
)

type SampleMessage struct {
	Name  string
	Value string
}

func main() {
	ctx := gctx.New()
	req, err := http.NewRequest("POST", "http://localhost:8080/stream", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	err = readBodyStream(ctx, resp.Body)
	if err != nil {
		panic(err)
	}
}

func readBodyStream(ctx context.Context, r io.Reader) error {
	scanner := streamio.NewScanner(r)
	for scanner.Scan() {
		messageContent := trimDataMessage(scanner.Text())
		g.Log().Debugf(ctx, "Message content: %s", messageContent)
		message := (*SampleMessage)(nil)
		if err := gjson.DecodeTo(messageContent, &message); err != nil {
			return err
		}
		g.Log().Infof(ctx, "Message decoded: SampleMessage %s", gjson.MustEncodeString(message))
	}
	return nil
}

// trimDataMessage strips possible "data:" prefix from a string, as well as
// whitespaces from the beginning and end of it.
func trimDataMessage(content string) string {
	if gstr.Pos(content, "data:") == 0 {
		return gstr.Trim(content[5:])
	}
	return gstr.Trim(content)
}
