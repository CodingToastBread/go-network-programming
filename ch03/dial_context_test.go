package ch03

import (
	"context"
	"errors"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	// deadline 정의
	dl := time.Now().Add(5 * time.Second)

	ctx, cancel := context.WithDeadline(context.Background(), dl)

	defer cancel()

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		_ = conn.Close()
		t.Fatal("connection did not time out")
	}

	var nErr net.Error
	if errors.As(err, &nErr) {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	} else {
		t.Error(err)
	}

	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("expected deadline exceeded; actual %v", ctx.Err())
	}
}
