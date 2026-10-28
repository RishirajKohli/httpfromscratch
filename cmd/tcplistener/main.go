package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(lines)
		data := make([]byte, 8)
		curLine := ""
		for {
			n, err := f.Read(data)
			if err != nil {
				fmt.Println("Error while reading happened", err)
				if err == io.EOF {
					lines <- curLine
					break
				}

				log.Fatal(err)
			}
			parts := strings.Split(string(data[:n]), "\n")
			if len(parts) == 1 {
				curLine += parts[0]
				continue
			}
			lines <- curLine + parts[0]
			for i := 1; i < len(parts)-1; i++ {
				lines <- parts[i]
			}
			curLine = parts[len(parts)-1]
		}
	}()

	return lines
}

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	fmt.Println("Listener setup complete. Awaiting connections")
	for {
		conn, err := listener.Accept()
		fmt.Println("Connection accepted")
		if err != nil {
			fmt.Println(err)
			break

		}
		defer conn.Close()

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection has been closed")
	}

}
