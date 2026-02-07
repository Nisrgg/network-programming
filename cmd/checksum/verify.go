package main

import (
	"crypto/sha512"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verify <downloaded_file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	if len(data) < 32 {
		panic("file too small to contain checksum")
	}

	payload := data[:len(data)-32]
	receivedSum := data[len(data)-32:]

	calculatedSum := sha512.Sum512_256(payload)

	if string(receivedSum) == string(calculatedSum[:]) {
		fmt.Println("OK: checksum matches")
	} else {
		fmt.Println("FAIL: checksum mismatch")
	}
}
