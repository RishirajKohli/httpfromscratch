package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 1. Resolve the UDP address for localhost:42069
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address: %v\n", err)
		os.Exit(1)
	}

	// 2. Prepare a UDP connection using net.DialUDP and defer closing it
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 3. Create a new bufio.Reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	// 4. Infinite loop to prompt, read, and send packets
	for {
		// Print prompt
		fmt.Print("> ")

		// Read a line from stdin
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			continue
		}

		// Write the line to the UDP connection
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to UDP connection: %v\n", err)
		}
	}
}
