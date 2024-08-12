package ntfy

import (
	"context"
	"os"
	"testing"
)

func TestSendNotify(t *testing.T) {
	host := os.Getenv("NTFY_HOST")
	token := os.Getenv("NTFY_TOKEN")
	topic := os.Getenv("NTFY_TOPIC")
	if host == "" || token == "" || topic == "" {
		t.Errorf("Notify env is unsufficient.")
		return
	}
	ctx := context.Background()
	cli := New(host, token)
	err := cli.Send(ctx, topic, HighLevel, "Title", "Body")
	if err != nil {
		t.Errorf("send notify message failed: %v", err)
		return
	}
}
