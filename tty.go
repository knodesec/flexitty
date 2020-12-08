package flexitty

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	historybuf "github.com/knodesec/go-historybuffer"
)

type TTY struct {
	Command string
	Args    []string
	History historybuf
	cmd     *exec.Cmd
	PTY     *os.File
}

func New(command string, argv []string) (*TTY, error) {

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
		History: historybuf.New(4096),
	}

	return newTTY, nil
}

func (t *TTY) Write(data []byte) error {
	_, err := t.PTY.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTY) Read() ([]byte, error) {
	data := make([]byte, 512)
	_, err := t.PTY.Read(data)
	if err != nil {
		return nil, err
	}
	_, err := t.History.Write(data)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func (t *TTY) Close() {
	t.PTY.Close()
}
