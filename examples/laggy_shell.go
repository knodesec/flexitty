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
	defer tty.Close() // Best effort.

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
	for data := range tty.OutputChan {
		for _, b := range data {
			if b > 0x20 && b < 0x7E {
				os.Stdout.Write([]byte{b})
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(MaxTypeDelay-MinTypeDelay+5)+MinTypeDelay))
				continue
			}
			os.Stdout.Write([]byte{b})
		}
	}
	log.Printf("Closing.\n")

	return nil
}

func main() {
	if err := test(); err != nil {
		log.Fatal(err)
	}
}
