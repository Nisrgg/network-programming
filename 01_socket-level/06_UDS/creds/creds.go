package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	auth "github.com/Nisrgg/network-programming/01_socket-level/06_UDS/creds/peer_authentication"
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n\t%s 1<group names>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func parseGroupNames(args []string) map[string]struct{} {
	groups := make(map[string]struct{})

	for _, arg := range args {
		grp, err := user.LookupGroup(arg)
		if err != nil {
			log.Println(err)
			continue
		}

		groups[grp.Gid] = struct{}{}
	}

	return groups
}

func main() {
	flag.Parse()
	groups := parseGroupNames(flag.Args())
	socket := filepath.Join(os.TempDir(), "creds.sock")
	addr, err := net.ResolveUnixAddr("unix", socket)

	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		_ = s.Close()
	}()

	fmt.Printf("Listening on %s ...\n", socket)
	for {
		conn, err := s.AcceptUnix()
		if err != nil {
			break
		}

		if auth.Allowed(conn, groups) {
			_, err = conn.Write([]byte("Welcome\n"))
			if err == nil {
				// handle the connection in a goroutine here
				continue
			}
		} else {
			_, err = conn.Write([]byte("Access denied\n"))
		}

		if err != nil {
			log.Println(err)
		}

		_ = conn.Close()
	}
}
