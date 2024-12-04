package main

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
