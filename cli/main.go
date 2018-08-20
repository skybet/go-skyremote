package main

import (
	"log"

	skyremote "github.com/skybet/go-skyremote"
	flag "github.com/spf13/pflag"
)

func main() {
	ip := flag.String("ip", "", "IP of remote box")
	port := flag.Int("port", 49160, "Port on remote box")
	channel := flag.String("channel", "100", "3 digit channel number")
	flag.Parse()
	if len(*ip) == 0 {
		log.Fatalf("-ip flag is required")
	}
	if len(*channel) != 3 {
		log.Fatalf("Channel must be 3 digits long")
	}
	s := skyremote.New(*ip, *port)
	if err := s.ChangeChannel(*channel); err != nil {
		log.Fatalf("Error changing channel: %s", err)
	}
}
