package call

import (
	"context"
	"crypto/sha256"
)

// MethodKey identifies a particular method on a component (formed by
// fingerprinting the component and method name).
type MethodKey [16]byte

// MakeMethodKey returns the fingerprint for the specified method on component.
func MakeMethodKey(component, method string) MethodKey {
	sig := sha256.Sum256([]byte(component + "." + method))
	var fp MethodKey
	copy(fp[:], sig[:])
	return fp
}

// Handler is a function that handles remote procedure calls. Regular
// application errors should be serialized in the returned bytes. A Handler
// should only return a non-nil error if the handler was not able to execute
// successfully.
type Handler func(ctx context.Context, args []byte) ([]byte, error)

// HandlerMap is a mapping from MethodID to a Handler. The zero value for a
// HandlerMap is an empty map.
type HandlerMap struct {
	handlers map[MethodKey]Handler
	names    map[MethodKey]string
}

// Set registers a handler for the specified method of component.
func (hm *HandlerMap) Set(component, method string, handler Handler) {
	if hm.handlers == nil {
		hm.handlers = map[MethodKey]Handler{}
		hm.names = map[MethodKey]string{}
	}
	fp := MakeMethodKey(component, method)
	hm.handlers[fp] = handler
	hm.names[fp] = component + "." + method
}
