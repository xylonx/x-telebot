package lolicon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/hystrix"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

const Endpoint = "https://api.lolicon.app/setu/v2"

type LoliconClient struct {
	httpClient *hystrix.Client
}

type ReqParameters struct {
	R18    bool     `json:"r18,omitempty"`
	Number int      `json:"num,omitempty"`
	UID    []int    `json:"uid,omitempty"`
	Tag    []string `json:"tag,omitempty"`
	Size   []string `json:"size,omitempty"`
	Proxy  string   `json:"proxy,omitempty"`
}

type Response struct {
	Error string `json:"error"`
	Data  []struct {
		PID    int      `json:"pid"`
		Page   int      `json:"p"`
		UID    int      `json:"uid"`
		Title  string   `json:"title"`
		Author string   `json:"author"`
		R18    bool     `json:"r18"`
		Tags   []string `json:"tags"`
		URLs   struct {
			Original string `json:"original"`
		} `json:"urls"`
	} `json:"data"`
}

var DefaultLoliconClient = &LoliconClient{
	httpClient: hystrix.NewClient(
		hystrix.WithHTTPTimeout(time.Second*2),
		hystrix.WithCommandName("lolicon_api"),
		hystrix.WithHystrixTimeout(time.Second*2),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(20),
		hystrix.WithRetryCount(3),
	),
}

func (c *LoliconClient) GetPictures(ctx context.Context, porn bool, num int, tags []string) (*Response, error) {
	data := ReqParameters{
		R18:    porn,
		Number: num,
		UID:    nil,
		Tag:    tags,
		Size:   []string{"original"},
		Proxy:  "https://i.pixiv.cat/{{path}}",
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		zapx.Error("encode request param failed", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, Endpoint, buf)
	if err != nil {
		zapx.Error("create request failed", zap.Error(err))
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		zapx.Error("http request failed", zap.Error(err))
		return nil, err
	}

	pics := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(pics); err != nil {
		zapx.Error("decode response failed", zap.Error(err))
		return nil, err
	}

	if pics.Error != "" {
		zapx.Error(pics.Error)
		return nil, errors.New(pics.Error)
	}

	return pics, nil
}
