package codegen

import "errors"

// CatchPanics recovers from panic() calls that occur during encoding,
// decoding, and RPC execution.
func CatchPanics(r interface{}) error {
	if r == nil {
		return nil
	}
	err, ok := r.(error)
	if !ok {
		panic(r)
	}
	if errors.As(err, &encoderError{}) || errors.As(err, &decoderError{}) {
		return err
	}
	panic(r)
}
