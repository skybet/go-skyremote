package main

import (
	"log"
	"strconv"

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
	s := skyremote.SkyRemote{Host: *ip, Port: *port}

	for _, n := range *channel {
		i, err := strconv.ParseInt(string(n), 10, 0)
		if err != nil {
			log.Fatalf("Error parsing integer: %s", err)
		}
		c, err := s.CommandFromDigit(int(i))
		if err != nil {
			log.Fatalf("Error parsing channel: %s", err)
		}
		if err := s.SendCommand(c); err != nil {
			log.Fatalf("Error sending command: %s", err)
		}
	}
}
