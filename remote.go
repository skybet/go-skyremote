package skyremote

import (
	"errors"
	"fmt"
	"math"
	"net"
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

// SkyRemote allows us to control Sky+HD and SkyQ boxes over IP
type SkyRemote struct {
	Host string
	Port int
}

// CommandFromDigit returns the correct command from a signle digit 0-9
func (s *SkyRemote) CommandFromDigit(n int) (Command, error) {
	if n > 9 {
		return Cmd0, errors.New("Must be a single digit between 0 and 9")
	}
	return Command(n + int(Cmd0)), nil
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
				// We also have to change the second byte to 0x00 and send again
				// for reasons unknown
				cmd[1] = 0
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
