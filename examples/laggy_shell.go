package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/knodesec/flexitty"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	MinTypeDelay = 10
	MaxTypeDelay = 30
)

func test() error {

	rand.Seed(time.Now().UnixNano())

	// Start the command with a pty.
	tty, err := flexitty.New("zsh", []string{})
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = tty.PTY.Close() }() // Best effort.

	// Fixed PTY size for demo
	tty.Resize(35, 50)

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() {
		for {
			data := make([]byte, 64)
			_, err := os.Stdin.Read(data)
			if err == io.EOF {
				log.Printf("Stdin read got EOF")
				break
			}
			tty.InputChan <- data
		}
	}()
	log.Printf("Starting chan range\n")
	for data := range tty.OutputChan {
		//log.Printf("Got data...\n")
		for _, b := range data {
			if b > 0x20 && b < 127 {
				os.Stdout.Write([]byte{b})
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(MaxTypeDelay-MinTypeDelay+5)+MinTypeDelay))
				continue
			}
			os.Stdout.Write([]byte{b})
		}
		//_, _ = os.Stdout.Write(data)
		//time.Sleep(time.Second * 2)
	}
	log.Printf("Closing.\n")

	//go func() { _, _ = io.Copy(tty.PTY, os.Stdin) }()
	//_, _ = io.Copy(os.Stdout, tty.PTY)

	return nil
}

func main() {
	if err := test(); err != nil {
		log.Fatal(err)
	}
}
