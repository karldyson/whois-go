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
	"runtime"
)

// define our variables
var (
	host		= flag.String("host", "whois.iana.org", "the whois server hostname")
	port		= flag.String("port", "43", "the whois service port to use")
	debugOutput	= flag.Bool("debug", false, "enable debug output")
	version		= flag.Bool("version", false, "the code version")
	revision	= flag.Bool("revision", false, "revision and build information")
	timeout 	= flag.Float64("timeout", 2, "the connection timeout value in seconds")
	versionString	string = "devel"
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
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(*host, *port), time.Duration(*timeout * float64(time.Second)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to %s failed: %s\n", net.JoinHostPort(*host, *port), err)
		return
	} else {
		_debug(fmt.Sprintf("Connection established to %s", net.JoinHostPort(*host, *port)))
	}

	// write to the server
	_, err = conn.Write([]byte(name + "\r\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Send to %s:%s failed: %s\n", *host, *port, err)
		return
	} else {
		_debug(fmt.Sprintf("data string [%s] sent to server", name))
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

func _debug(debugString string) {
	if !*debugOutput {
		return
	}
	pc, _, no, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fmt.Printf("DEBUG :: %s#%d :: %s\n", details.Name(), no, debugString)
	} else {
		fmt.Fprintf(os.Stderr, "Error: fatal error determining debug calling function")
		os.Exit(1)
	}
}
