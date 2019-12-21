package flexitty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
)

type TTY struct {
	Command string
	Args    []string

	InputChan  chan []byte
	OutputChan chan []byte
	ptyClosed  <-chan int

	cmd *exec.Cmd
	PTY *os.File
}

func New(command string, argv []string) (*TTY, error) {

	fmt.Printf("Starting new FlexiTTY with %s\n", command)
	cmd := exec.Command(command, argv...)

	pty, err := pty.Start(cmd)
	if err != nil {
		// todo close cmd?
		return nil, fmt.Errorf("failed to start pty: %s", err.Error)
	}

	newTTY := &TTY{
		Command: command,
		Args:    argv,
		cmd:     cmd,
		PTY:     pty,
	}

	newTTY.OutputChan = make(chan []byte)
	newTTY.InputChan = make(chan []byte)
	newTTY.StartChannels()

	return newTTY, nil
}

func (t *TTY) Write(data []byte) error {
	_, err := t.PTY.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTY) Read(data []byte) error {
	_, err := t.PTY.Read(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTY) Resize(width, height int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		t.PTY.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func (t *TTY) StartChannels() {
	go func() {
		//log.Printf("Started go subroutine channel\n")
		for {
			data := make([]byte, 512)
			err := t.Read(data)
			if err == io.EOF {
				t.Close()
				return
			} else if err != nil {
				t.Close()
				return
			}
			//log.Printf("Outputting to channel\n")
			t.OutputChan <- data
		}
	}()

	go func() {
		for data := range t.InputChan {
			err := t.Write(data)
			if err != nil {
				t.Close()
				panic(err)
			}
		}
	}()
}

func (t *TTY) Close() {
	close(t.OutputChan)
	close(t.InputChan)
	t.PTY.Close()
}
