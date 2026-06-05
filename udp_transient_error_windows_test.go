//go:build windows

package raknet

import (
	"fmt"
	"syscall"
	"testing"
)

func TestIsTransientUDPReadErrorIncludesMessageTooLong(t *testing.T) {
	err := fmt.Errorf("read udp: %w", syscall.Errno(10040))
	if !isTransientUDPReadError(err) {
		t.Fatal("expected WSAEMSGSIZE to be treated as a transient UDP read error")
	}
}
