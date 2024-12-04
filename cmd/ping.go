package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	icmpEchoRequest = 8
	icmpEchoReply   = 0
)

func sendPing(addr string, id, seq uint16) (time.Duration, error) {
	conn, err := net.Dial("ip4:icmp", addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	packet := icmpPacket{
		Type: icmpEchoRequest,
		Code: 0,
		ID:   id,
		Seq:  seq,
	}
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint16(buffer[4:], id)
	binary.BigEndian.PutUint16(buffer[6:], seq)
	packet.Checksum = checksum(buffer)
	binary.BigEndian.PutUint16(buffer[2:], packet.Checksum)

	start := time.Now()
	if _, err := conn.Write(buffer); err != nil {
		return 0, err
	}

	reply := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Read(reply)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	if reply[20] != icmpEchoReply {
		return 0, fmt.Errorf("did not receive echo reply")
	}

	return elapsed, nil
}
