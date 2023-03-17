package logging

import (
	"errors"
	"io"
)

// tailReader is an io.Reader that behaves like tail -f.
type tailReader struct {
	src            io.Reader    // the reader from which we read
	waitForChanges func() error // blocks until src has bytes to read
}

var _ io.Reader = &tailReader{}

// newTailReader returns a reader that yields data from src. On end-of-file, it
// waits for the reader to grow by calling waitForChanges instead of returning
// io.EOF (and exits immediately if waitForChanges returns an error).
func newTailReader(src io.Reader, waitForChanges func() error) *tailReader {
	return &tailReader{
		src:            src,
		waitForChanges: waitForChanges,
	}
}

// Read returns available data, or waits for more data to become available.
// Read implements the io.Reader interface.
func (t *tailReader) Read(p []byte) (int, error) {
	for {
		n, err := t.src.Read(p)
		if !errors.Is(err, io.EOF) {
			return n, err
		}

		if n > 0 {
			// Return what we got.
			return n, nil
		}

		if err := t.waitForChanges(); err != nil {
			return 0, err
		}
	}
}
