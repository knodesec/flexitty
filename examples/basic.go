package main

import (
	"io"
	"log"
	"os"

	"github.com/knodesec/flexitty"
)

func main() {

	log.SetFlags(0)

	log.Printf("Basic FlexiTTY Usage\n")

	tty, err := flexitty.New("ls", []string{"-l"})
	if err != nil {
		log.Fatal(err)
	}

	/*err = tty.Write([]byte("TestString!"))
	if err != nil {
		log.Fatal(err)
	}
	*/
	io.Copy(os.Stdout, tty.PTY)

}
