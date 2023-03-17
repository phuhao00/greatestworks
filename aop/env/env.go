// Package env implements helper functions for dealing with environment variables.
package env

import (
	"fmt"
	"strings"
)

// Split splits an environment variable of the form key=value into its
// constituent key and value parts.
func Split(kv string) (string, string, error) {
	k, v, ok := strings.Cut(kv, "=")
	if !ok {
		return "", "", fmt.Errorf("env: %q is not of form key=value", kv)
	}
	return k, v, nil
}

// Parse parses a list of environment variables of the form key=value into a
// map from keys to values. If a key appears multiple times, the last value of
// the key is returned.
func Parse(env []string) (map[string]string, error) {
	kvs := map[string]string{}
	for _, kv := range env {
		k, v, err := Split(kv)
		if err != nil {
			return nil, err
		}
		kvs[k] = v
	}
	return kvs, nil
}
