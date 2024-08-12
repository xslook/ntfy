package ntfy

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	host  string
	token string
}

func New(host, token string) *Client {
	return &Client{
		host:  host,
		token: token,
	}
}

type Level int

func (r Level) String() string {
	return strconv.Itoa(int(r))
}

const (
	UnknownLevel Level = iota
	MinLevel
	LowLevel
	NormalLevel
	HighLevel
	MaxLevel
)

type Message struct {
	Title    string
	Body     string
	Priority Level
	Tags     []string
	Attach   string
}

func (cli *Client) Send(ctx context.Context, topic string, p Level, title, body string) error {
	return cli.SendMessage(ctx, topic, &Message{
		Priority: p,
		Title:    title,
		Body:     body,
	})
}

func (cli *Client) SendMessage(ctx context.Context, topic string, msg *Message) error {
	if msg == nil {
		return nil
	}
	if msg.Body == "" {
		return fmt.Errorf("invalid empty message body")
	}
	addr, err := url.JoinPath(cli.host, topic)
	if err != nil {
		return fmt.Errorf("Join url path failed: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, strings.NewReader(msg.Body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cli.token)
	if msg.Title != "" {
		req.Header.Set("X-Title", msg.Title)
	}
	if msg.Priority > UnknownLevel {
		req.Header.Set("X-Priority", msg.Priority.String())
	}
	if len(msg.Tags) > 0 {
		req.Header.Set("X-Tags", strings.Join(msg.Tags, ","))
	}
	if msg.Attach != "" {
		req.Header.Set("X-Attach", msg.Attach)
	}
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("do http request failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http response not ok, %s", res.Status)
	}
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	slog.With("response", string(bts)).DebugContext(ctx, "Ntfy send response")
	return nil
}
