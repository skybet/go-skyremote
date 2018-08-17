package main

import (
	"fmt"
	"log"
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
func (s *SkyRemote) SendCommand(c Command) error {
	// Create the byte sequence
	cmd := []byte{4, 1, 0, 0, 0, 0, byte(int(math.Floor(224 + float64((c / 16))))), byte(c % 16)}

	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	//defer conn.Close()
	if err != nil {
		return err
	}
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return err
	}
	log.Printf("Received: % x\n", buff[:n])

	log.Printf("Sending: % x\n", buff[:12])
	n, err = conn.Write(buff[:12])
	if err != nil {
		return err
	}
	n, err = conn.Read(buff)
	if err != nil {
		return err
	}
	log.Printf("Received: % x\n", buff[:n])

	log.Printf("Sending: % x\n", buff[:1])
	n, err = conn.Write(buff[:1])
	if err != nil {
		return err
	}
	n, err = conn.Read(buff)
	if err != nil {
		return err
	}
	log.Printf("Received: % x\n", buff[:n])

	log.Printf("Sending: % x\n", buff[:1])
	n, err = conn.Write(buff[:1])
	if err != nil {
		return err
	}
	n, err = conn.Read(buff)
	if err != nil {
		return err
	}
	log.Printf("Received: % x\n", buff[:n])

	log.Printf("Sending: % x\n", cmd)
	n, err = conn.Write(cmd)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func main() {
	s := SkyRemote{Host: "172.16.4.51", Port: 49160}
	if err := s.SendCommand(CmdGreen); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
