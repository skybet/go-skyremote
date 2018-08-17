package skyremote

import (
	"fmt"
	"math"
	"net"
)

// Command represents a button on the remote
type Command int

// Set of commands (buttons) which can be sent
const (
	CmdPower       = 0
	CmdSelect      = 1
	CmdBackup      = 2
	CmdDismiss     = 2
	CmdChannelup   = 6
	CmdChanneldown = 7
	CmdInteractive = 8
	CmdSidebar     = 8
	CmdHelp        = 9
	CmdServices    = 10
	CmdSearch      = 10
	CmdTvguide     = 11
	CmdHome        = 11
	CmdI           = 14
	CmdText        = 15
	CmdUp          = 16
	CmdDown        = 17
	CmdLeft        = 18
	CmdRight       = 19
	CmdRed         = 32
	CmdGreen       = 33
	CmdYellow      = 34
	CmdBlue        = 35
	Cmd0           = 48
	Cmd1           = 49
	Cmd2           = 50
	Cmd3           = 51
	Cmd4           = 52
	Cmd5           = 53
	Cmd6           = 54
	Cmd7           = 55
	Cmd8           = 56
	Cmd9           = 57
	CmdPlay        = 64
	CmdPause       = 65
	CmdStop        = 66
	CmdRecord      = 67
	CmdFastforward = 69
	CmdRewind      = 71
	CmdBoxoffice   = 240
	CmdSky         = 241
)

// SkyRemote allows us to control Sky+HD and SkyQ boxes over IP
type SkyRemote struct {
	Host string
	Port int
}

// SendCommand to Sky box
func (s *SkyRemote) SendCommand(c Command) (err error) {
	// Attempt to connect
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}

	// This will hold data from the socket while we examine it
	buff := make([]byte, 256)
	// Create the byte sequence
	cmd := []byte{4, 1, 0, 0, 0, 0, byte(int(math.Floor(224 + float64((c / 16))))), byte(c % 16)}
	// When this is true, we are done
	var done bool

	for !done {
		if n, err := conn.Read(buff); err == nil {
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
				done = true
			// Dunno - panic?
			default:
				err = fmt.Errorf("Unexpected byte sequence received: % x")
				done = true
			}
		} else {
			break
		}
	}
	// Clean up
	conn.Close()
	return
}
