package packer

import (
	"errors"
	"io"
	"log"
	"strings"
)

// TestUi is a UI that only outputs unformatted output to the given writer.
type TestUi struct {
	Writer io.Writer
}

var _ Ui = new(TestUi)

func (u *TestUi) Ask(query string) (string, error) {
	return "", errors.New("Test UI can't ask")
}

func (u *TestUi) Say(message string) {
	log.Printf(message)
}

func (u *TestUi) Message(message string) {
	log.Printf(message)
}

func (u *TestUi) Error(message string) {
	log.Printf(message)
}

func (u *TestUi) Machine(message string, args ...string) {
	log.Printf("%s, %s", message, strings.Join(args, ", "))
}

func (u *TestUi) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
	return &readCloser{
		read: func(p []byte) (int, error) {
			return stream.Read(p)
		},
		close: func() error {
			return stream.Close()
		},
	}
}

type MockUi struct {
	AskCalled      bool
	AskQuery       string
	ErrorCalled    bool
	ErrorMessage   string
	MachineCalled  bool
	MachineType    string
	MachineArgs    []string
	MessageCalled  bool
	MessageMessage string
	SayCalled      bool
	SayMessage     string

	TrackProgressCalled    bool
	ProgressBarAddCalled   bool
	ProgressBarCloseCalled bool
}

func (u *MockUi) Ask(query string) (string, error) {
	u.AskCalled = true
	u.AskQuery = query
	return "foo", nil
}

func (u *MockUi) Error(message string) {
	u.ErrorCalled = true
	u.ErrorMessage = message
}

func (u *MockUi) Machine(t string, args ...string) {
	u.MachineCalled = true
	u.MachineType = t
	u.MachineArgs = args
}

func (u *MockUi) Message(message string) {
	u.MessageCalled = true
	u.MessageMessage = message
}

func (u *MockUi) Say(message string) {
	u.SayCalled = true
	u.SayMessage = message
}

func (u *MockUi) TrackProgress(_ string, _, _ int64, stream io.ReadCloser) (body io.ReadCloser) {
	u.TrackProgressCalled = true

	return &readCloser{
		read: func(p []byte) (int, error) {
			u.ProgressBarAddCalled = true
			return stream.Read(p)
		},
		close: func() error {
			u.ProgressBarCloseCalled = true
			return stream.Close()
		},
	}
}

type readCloser struct {
	read  func([]byte) (int, error)
	close func() error
}

func (c *readCloser) Close() error               { return c.close() }
func (c *readCloser) Read(p []byte) (int, error) { return c.read(p) }
