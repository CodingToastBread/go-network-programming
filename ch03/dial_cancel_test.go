package ch03

import (
	"context"
	"errors"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sync := make(chan struct{})

	go func() {
		defer func() { sync <- struct{}{} }()

		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err) // 여기에 오면 사실상 성공임.
			return
		}
		_ = conn.Close()
		t.Error("connection did not time out")
	}()

	cancel()
	<-sync

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Errorf("expected canceled context; actual %q", ctx.Err())
	}
}
