package raknet_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Yeah114/go-raknet"
)

func TestListen(t *testing.T) {
	l, err := raknet.Listen("127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	dialErr := make(chan error, 1)
	go func() {
		conn, err := raknet.Dial(l.Addr().String())
		if conn != nil {
			_ = conn.Close()
		}
		dialErr <- err
	}()
	c := make(chan error)
	go accept(l, c)

	select {
	case err := <-c:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("accepting connection took longer than 3 seconds")
	}
	select {
	case err := <-dialErr:
		if err != nil {
			t.Errorf("error dialing connection: %v", err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("dialing connection took longer than 3 seconds")
	}
}

func TestListenVersion(t *testing.T) {
	const protocolVersion byte = 12

	l, err := raknet.ListenVersion("127.0.0.1:0", protocolVersion)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	dialErr := make(chan error, 1)
	go func() {
		conn, err := raknet.DialVersion(l.Addr().String(), protocolVersion)
		if conn != nil {
			_ = conn.Close()
		}
		dialErr <- err
	}()
	c := make(chan error)
	go accept(l, c)

	select {
	case err := <-c:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("accepting connection took longer than 3 seconds")
	}
	select {
	case err := <-dialErr:
		if err != nil {
			t.Errorf("error dialing connection: %v", err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("dialing connection took longer than 3 seconds")
	}
}

func TestListenConfigProtocolVersion(t *testing.T) {
	const protocolVersion byte = 12

	l, err := raknet.ListenConfig{ProtocolVersion: protocolVersion}.Listen("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	dialer := raknet.Dialer{ProtocolVersion: protocolVersion}
	dialErr := make(chan error, 1)
	go func() {
		conn, err := dialer.Dial(l.Addr().String())
		if conn != nil {
			_ = conn.Close()
		}
		dialErr <- err
	}()
	c := make(chan error)
	go accept(l, c)

	select {
	case err := <-c:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("accepting connection took longer than 3 seconds")
	}
	select {
	case err := <-dialErr:
		if err != nil {
			t.Errorf("error dialing connection: %v", err)
		}
	case <-time.After(time.Second * 3):
		t.Errorf("dialing connection took longer than 3 seconds")
	}
}

func TestListenAutoVersion(t *testing.T) {
	l, err := raknet.ListenAutoVersion("127.0.0.1:0", 11, 12)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for _, protocolVersion := range []byte{11, 12} {
		dialErr := make(chan error, 1)
		go func() {
			conn, err := raknet.DialVersion(l.Addr().String(), protocolVersion)
			if conn != nil {
				_ = conn.Close()
			}
			dialErr <- err
		}()
		c := make(chan error)
		go accept(l, c)

		select {
		case err := <-c:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(time.Second * 3):
			t.Errorf("accepting version %v connection took longer than 3 seconds", protocolVersion)
		}
		select {
		case err := <-dialErr:
			if err != nil {
				t.Errorf("error dialing version %v connection: %v", protocolVersion, err)
			}
		case <-time.After(time.Second * 3):
			t.Errorf("dialing version %v connection took longer than 3 seconds", protocolVersion)
		}
	}

	if conn, err := raknet.DialTimeoutVersion(l.Addr().String(), time.Second, 13); err == nil {
		_ = conn.Close()
		t.Fatal("expected mismatched protocol error")
	}
}

func TestListenVersionRejectsMismatchedDialVersion(t *testing.T) {
	l, err := raknet.ListenVersion("127.0.0.1:0", 12)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	if conn, err := raknet.DialTimeoutVersion(l.Addr().String(), time.Second, 13); err == nil {
		_ = conn.Close()
		t.Fatal("expected mismatched protocol error")
	}
}

func accept(l *raknet.Listener, c chan error) {
	if _, err := l.Accept(); err != nil {
		c <- fmt.Errorf("error accepting connection: %v", err)
	}
	c <- nil
}
