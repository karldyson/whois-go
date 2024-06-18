package main

import (
	"fmt"
	"net"
	"flag"
	"os"
	"bufio"
	"io"
	"time"
)

// define our variables
var (
        host = flag.String("host", "", "the whois hostname")
	port = flag.String("port", "43", "the port to use")
	domain = flag.String("domain", "", "the domain name to query")
)

func main() {

	// define the usage output
	flag.Usage = func() {
                fmt.Fprintf(os.Stderr, "Usage: %s --host <host --domain <domain> [--port <port>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
                flag.PrintDefaults()
	}

	// parse the flags
	flag.Parse()

	// bail if we don't have what we need
	if *host == "" || *domain == "" {
		flag.Usage()
		return
	}

	// connect to the server
	m, _ := time.ParseDuration("2s")
	dialer := net.Dialer{Timeout: m}
	conn, err := dialer.Dial("tcp", *host + ":" + *port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to %s:%s failed: %s\n", *host, *port, err)
		return
	}

	// write to the server
	_, err = conn.Write([]byte(*domain + "\r\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Send to %s:%s failed: %s\n", *host, *port, err)
		return
	}
	
	// receive the response from the server
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes(byte('\n'))
		switch err {
		case nil:
			break
		case io.EOF:
			return
		default:
			fmt.Fprintf(os.Stderr, "Receive from %s:%s failed: %s\n", *host, *port, err)
			return
		}
		fmt.Printf("%s", line)
	}	
	
	// close the connection
	conn.Close()
}
