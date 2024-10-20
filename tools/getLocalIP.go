package tools

import (
	"fmt"
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetLocalIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return fmt.Sprint(localAddr.IP)
}
