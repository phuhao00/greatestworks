package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"greatestworks/aop/protos"
)

// exampleReader is a toy Reader implementation used by ExampleReader.
type exampleReader struct {
	entries []*protos.LogEntry
	closed  bool
}

// Read implements the Reader interface.
func (r *exampleReader) Read(context.Context) (*protos.LogEntry, error) {
	if r.closed {
		return nil, fmt.Errorf("closed")
	}
	if len(r.entries) == 0 {
		return nil, io.EOF
	}
	entry := r.entries[0]
	r.entries = r.entries[1:]
	return entry, nil
}

// Close implements the Reader interface.
func (r *exampleReader) Close() {
	r.closed = true
}

func getLogReader() *exampleReader {
	return &exampleReader{
		entries: []*protos.LogEntry{
			{Msg: "1"},
			{Msg: "2"},
			{Msg: "3"},
		},
	}
}

func ExampleReader() {
	ctx := context.Background()
	reader := getLogReader()
	defer reader.Close()
	for {
		entry, err := reader.Read(ctx)
		if errors.Is(err, io.EOF) {
			// No more log entries.
			return
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println(entry)
	}
}

// TestDontShowWholeFile is here so that this entire file isn't shown as an
// example. As explained in [1],
//
// > The entire test file is presented as the example when it contains a single
// > example function, at least one other function, type, variable, or constant
// > declaration, and no test or benchmark functions.
//
// [1]: https://pkg.go.dev/testing#hdr-Examples
func TestDontShowWholeFile(*testing.T) {}
