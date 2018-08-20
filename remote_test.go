package skyremote

import (
	"bytes"
	"net"
	"sync"
	"testing"
)

type stanza struct {
	request, response []byte
}

func TestSendCommand(t *testing.T) {
	// Fake a client and server using a pipe
	server, client := net.Pipe()
	defer client.Close()
	defer server.Close()

	testDialer := func(host string, port int) (net.Conn, error) {
		return client, nil
	}

	var wg sync.WaitGroup
	var e error
	go func() {
		wg.Add(1)
		s := SkyRemote{Dialer: testDialer}
		e = s.SendCommand(CmdTvguide)
		wg.Done()
	}()
	handshake := []stanza{
		stanza{
			request:  []byte{0x53, 0x4b, 0x59, 0x20, 0x30, 0x30, 0x30, 0x2e, 0x30, 0x30, 0x31, 0x0a},
			response: []byte{0x53, 0x4b, 0x59, 0x20, 0x30, 0x30, 0x30, 0x2e, 0x30, 0x30, 0x31, 0x0a},
		},
		stanza{
			request:  []byte{0x01, 0x01},
			response: []byte{0x01},
		},
		stanza{
			request:  []byte{0x00, 0x00, 0x00, 0x000},
			response: []byte{0x00},
		},
		stanza{
			request: []byte{0x00, 0x00, 0x00, 0x000, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x000, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			response: []byte{0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x0b},
		},
	}
	keyUp := []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x0b}
	buff := make([]byte, 256)

	for _, hs := range handshake {
		_, err := server.Write(hs.request)
		if err != nil {
			t.Fatalf("Error writing to pipe: %s", err)
		}
		n, err := server.Read(buff)
		if err != nil {
			t.Fatalf("Error reading from pipe: %s", err)
		}
		if bytes.Compare(hs.response, buff[:n]) != 0 {
			t.Fatal("Expected % x got % x", hs.response, buff[:n])
		}
	}
	n, err := server.Read(buff)
	if err != nil {
		t.Fatalf("Error reading from pipe: %s", err)
	}
	if bytes.Compare(keyUp, buff[:n]) != 0 {
		t.Fatal("Expected % x got % x", keyUp, buff[:n])
	}
	wg.Wait()
	if e != nil {
		t.Fatalf("Unexpected error from client: %s", e)
	}
}

func TestBadHandshake(t *testing.T) {
	server, client := net.Pipe()
	defer client.Close()
	defer server.Close()

	testDialer := func(host string, port int) (net.Conn, error) {
		return client, nil
	}
	var wg sync.WaitGroup
	var e error
	go func() {
		wg.Add(1)
		s := SkyRemote{Dialer: testDialer}
		e = s.SendCommand(CmdTvguide)
		wg.Done()
	}()
	hs := stanza{
		request:  []byte{0x53, 0x4b, 0x59, 0x20, 0x30, 0x30, 0x30, 0x2e, 0x30, 0x30, 0x31, 0x0a, 0x00},
		response: []byte{},
	}
	buff := make([]byte, 256)

	_, err := server.Write(hs.request)
	if err != nil {
		t.Fatalf("Error writing to pipe: %s", err)
	}
	_, err = server.Read(buff)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Error reading from pipe: %s", err)
	}
	wg.Wait()
	if e == nil {
		t.Fatal("This should be an error!")
	}
	if e.Error() != "Unexpected byte sequence received: 53 4b 59 20 30 30 30 2e 30 30 31 0a 00" {
		t.Fatalf("Unexpected error string: %s", e)
	}
}
