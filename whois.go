package main

import (
	"fmt"
	"net"
	"flag"
	"os"
	"bufio"
	"io"
	"time"
	"runtime/debug"
)

// define our variables
var (
	host		= flag.String("host", "whois.iana.org", "the whois server hostname")
	port		= flag.String("port", "43", "the whois service port to use")
	debugOutput	= flag.Bool("debug", false, "enable debug output")
	version		= flag.Bool("version", false, "the code version")
	revision	= flag.Bool("revision", false, "revision and build information")
	versionString string = "devel"
	name string
)

func main() {

	// define the usage output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --host <host> [--port <port>] <name>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	// parse the flags
	flag.Parse()

	// output revision info
	if *revision {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			panic("not ok reading build info!")
		}
		fmt.Printf("%s version information:\ncommit %s\n%+v\n", os.Args[0], versionString, bi)
		return
	}

	// output version info
	if *version {
		fmt.Printf("%s version %s\n", os.Args[0], versionString)
		return
	}

	switch len(flag.Args()) {
	case 0:
		fmt.Fprintf(os.Stderr, "ERROR: no string supplied\n")
		flag.Usage()
		os.Exit(1)
	case 1:
		name = flag.Args()[0]
	default:
	// and a default catching something weird
		fmt.Fprintf(os.Stderr, "ERROR: invalid number of CLI parameters\n")
		flag.Usage()
		os.Exit(1)
	}

	// connect to the server
	m, _ := time.ParseDuration("2s")
	dialer := net.Dialer{Timeout: m}
	conn, err := dialer.Dial("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to %s:%s failed: %s\n", *host, *port, err)
		return
	}

	// write to the server
	_, err = conn.Write([]byte(name + "\r\n"))
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
