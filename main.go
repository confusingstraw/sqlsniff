package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"net"
	"strconv"
)

var address = flag.String("address", "127.0.0.1", "Host address to connect to, without port")
var port = flag.Uint("port", 3306, "Port to connect to")

func init() {
	flag.Parse()
}

func main() {
	conn, err := net.Dial("tcp", *address+":"+strconv.FormatUint(uint64(*port), 10))
	if err != nil {
		log.Fatal(errors.New("failed to connect to server"))
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	p, n, err := ReadPacket(reader)
	if err != nil {
		log.Fatal(err)
	}

	desc, err := ParseDescriptor(p, n)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("AuthenticationPlugin: %s\n", desc.AuthenticationPlugin)
	log.Printf("Capabilities: %d\n", desc.Capabilities)
	log.Printf("ConnectionId: %d\n", desc.ConnectionId)
	log.Printf("ProtocolVersion: %d\n", desc.ProtocolVersion)
	log.Printf("ServerVersion: %s\n", desc.ServerVersion)
	log.Printf("Status: %d\n", desc.Status)
}
