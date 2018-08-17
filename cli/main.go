package main

import (
	"log"

	skyremote "github.com/skybet/go-skyremote"
	flag "github.com/spf13/pflag"
)

func main() {
	ip := flag.String("ip", "", "IP of remote box")
	port := flag.Int("port", 49160, "Port on remote box")
	flag.Parse()
	if len(*ip) == 0 {
		log.Fatalf("-ip flag is required")
	}
	s := skyremote.SkyRemote{Host: *ip, Port: *port}
	if err := s.SendCommand(skyremote.CmdTvguide); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
