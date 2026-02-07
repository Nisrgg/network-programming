package main

import (
	_ "crypto/sha512"
	"flag"
	"io/ioutil"
	"log"

	"github.com/Nisrgg/network-programming/01_socket-level/05_UDP/tftp"
)

var (
	address = flag.String("a", "127.0.0.1:69", "listen address")
	payload = flag.String("p", "cmd/payload.svg", "file to serve to clients")
)

func main() {
	flag.Parse()
	p, err := ioutil.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	// sum := sha512.Sum512_256(p)
	// p = append(p, sum[:]...) // append 32-byte checksum

	s := tftp.Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))
}
