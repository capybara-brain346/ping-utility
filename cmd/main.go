package main

import (
	"fmt"
	"os"
	"time"
)

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
