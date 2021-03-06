package skyremote

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
)

// Command represents a button on the remote
type Command int

// Set of commands (buttons) which can be sent
const (
	CmdPower       Command = 0
	CmdSelect      Command = 1
	CmdBackup      Command = 2
	CmdDismiss     Command = 2
	CmdChannelup   Command = 6
	CmdChanneldown Command = 7
	CmdInteractive Command = 8
	CmdSidebar     Command = 8
	CmdHelp        Command = 9
	CmdServices    Command = 10
	CmdSearch      Command = 10
	CmdTvguide     Command = 11
	CmdHome        Command = 11
	CmdI           Command = 14
	CmdText        Command = 15
	CmdUp          Command = 16
	CmdDown        Command = 17
	CmdLeft        Command = 18
	CmdRight       Command = 19
	CmdRed         Command = 32
	CmdGreen       Command = 33
	CmdYellow      Command = 34
	CmdBlue        Command = 35
	Cmd0           Command = 48
	Cmd1           Command = 49
	Cmd2           Command = 50
	Cmd3           Command = 51
	Cmd4           Command = 52
	Cmd5           Command = 53
	Cmd6           Command = 54
	Cmd7           Command = 55
	Cmd8           Command = 56
	Cmd9           Command = 57
	CmdPlay        Command = 64
	CmdPause       Command = 65
	CmdStop        Command = 66
	CmdRecord      Command = 67
	CmdFastforward Command = 69
	CmdRewind      Command = 71
	CmdBoxoffice   Command = 240
	CmdSky         Command = 241
)

// DialFunc is a function used to connect to remote boxes
type DialFunc func(host string, port int) (net.Conn, error)

// SkyRemote allows us to control Sky+HD and SkyQ boxes over IP
type SkyRemote struct {
	Host   string
	Port   int
	Dialer DialFunc
}

// New SkyRemote with default TCP dialer
func New(host string, port int) *SkyRemote {
	return &SkyRemote{
		Host:   host,
		Port:   port,
		Dialer: TCPDialer,
	}
}

// CommandFromDigit returns the correct command from a single digit 0-9
func (s *SkyRemote) CommandFromDigit(n int) (Command, error) {
	if n > 9 {
		return Cmd0, errors.New("Must be a single digit between 0 and 9")
	}
	return Command(n + int(Cmd0)), nil
}

// ChangeChannel on a Sky box, takes in a string e.g. "115"
func (s *SkyRemote) ChangeChannel(ch string) error {
	for _, n := range ch {
		i, err := strconv.ParseInt(string(n), 10, 0)
		if err != nil {
			return fmt.Errorf("invalid channel (%s), must be numberic", ch)
		}
		c, err := s.CommandFromDigit(int(i))
		if err != nil {
			return err
		}
		if err := s.SendCommand(c); err != nil {
			return err
		}
	}
	return nil
}

// SendCommand to Sky box
func (s *SkyRemote) SendCommand(c Command) error {
	// Attempt to connect
	conn, err := s.Dialer(s.Host, s.Port)
	if err != nil {
		return err
	}

	// This will hold data from the socket while we examine it
	buff := make([]byte, 256)
	// Create the byte sequence
	cmd := []byte{4, 1, 0, 0, 0, 0, byte(int(math.Floor(224 + float64((c / 16))))), byte(c % 16)}

loop:
	for err == nil {
		n, e := conn.Read(buff)
		if e != nil {
			return e
		}

		switch {
		// The first part of the handshake - 12 byte sequence received
		// we are to return it back verbatim
		case n == 12:
			_, err = conn.Write(buff[:12])
		// Two short sequences follow
		// we are to return the first byte of each
		case n < 12:
			_, err = conn.Write(buff[:1])
		// The final sequence received is 24 0x00s
		// after this we can send out command sequence
		case n == 24:
			_, err = conn.Write(cmd)
			// We also have to change the second byte to 0x00 and send again
			// for reasons unknown
			cmd[1] = 0
			_, err = conn.Write(cmd)
			break loop
		// Dunno - panic?
		default:
			err = fmt.Errorf("Unexpected byte sequence received: % x", buff[:n])
			break loop
		}
	}
	// Clean up
	conn.Close()
	return err
}

// TCPDialer returns a TCP connection
func TCPDialer(host string, port int) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout("tcp4", addr.String(), 2*time.Second)
	if err != nil {
		// handle error
		return nil, err
	}
	return conn, nil
}
