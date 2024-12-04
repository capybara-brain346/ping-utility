package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	icmpEchoRequest = 8
	icmpEchoReply   = 0
)

type icmpPacket struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
}

func checksum(data []byte) uint16 {
	sum := 0
	for i := 0; i < len(data)-1; i += 2 {
		sum += int(data[i])<<8 | int(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += int(data[len(data)-1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	return uint16(^sum)
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <IP>")
		return
	}

	addr := os.Args[1]
	fmt.Printf("PING %s:\n", addr)
	id := uint16(os.Getpid() & 0xffff)
	var sent, received int
	var rttTimes []time.Duration

	for i := 0; i < 4; i++ {
		rtt, err := sendPing(addr, id, uint16(i))
		sent++
		if err != nil {
			fmt.Printf("Request timed out: %v\n", err)
		} else {
			received++
			rttTimes = append(rttTimes, rtt)
			fmt.Printf("Reply from %s: time=%v\n", addr, rtt)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n--- Statistics ---")
	fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n", sent, received, float64(sent-received)/float64(sent)*100)
	if len(rttTimes) > 0 {
		var total time.Duration
		min, max := rttTimes[0], rttTimes[0]
		for _, rtt := range rttTimes {
			total += rtt
			if rtt < min {
				min = rtt
			}
			if rtt > max {
				max = rtt
			}
		}
		fmt.Printf("RTT: min=%v avg=%v max=%v\n", min, total/time.Duration(len(rttTimes)), max)
	}
}
