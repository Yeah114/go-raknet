//go:build !windows

package raknet

import (
	"fmt"
	"syscall"
	"testing"
)

func TestIsTransientUDPReadErrorIncludesMessageTooLong(t *testing.T) {
	err := fmt.Errorf("read udp: %w", syscall.EMSGSIZE)
	if !isTransientUDPReadError(err) {
		t.Fatal("expected EMSGSIZE to be treated as a transient UDP read error")
	}
}
