package versioned_map

import (
	"context"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
	"greatestworks/aop/cond"
	"greatestworks/aop/protomsg"
)

const missing = "__tombstone__"

type versioned[T proto.Message] struct {
	value   T
	version string
}

// Map is a simple, versioned_map, map.
type Map[T proto.Message] struct {
	mu     sync.Mutex
	global int
	data   map[string]versioned[T]

	// changed is a condition variable that wraps m. It's used to notify
	// versioned Gets when a value has changed.
	changed *cond.Cond
}

func NewMap[T proto.Message]() *Map[T] {
	s := Map[T]{global: 0, data: map[string]versioned[T]{}}
	s.changed = cond.NewCond(&s.mu)
	return &s
}

// Update (over)writes the value for a key with value.
func (m *Map[T]) Update(key string, value T) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	version := strconv.Itoa(m.global)
	m.global += 1
	m.data[key] = versioned[T]{protomsg.Clone(value), version}
	m.changed.Broadcast()
	return version
}

// Read gets the value and the version associated with a key.
//
// If version is the empty string, then read gets the latest value of the key,
// along with its version; otherwise it blocks until the latest value of the key
// is newer than the provided version.
func (m *Map[T]) Read(ctx context.Context, key string, version string) (T, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if version == missing {
		value, latest := m.getValue(key)
		return protomsg.Clone(value), latest, nil
	}
	for !m.hasChanged(key, version) {
		if err := m.changed.Wait(ctx); err != nil {
			var zero T
			return zero, missing, err
		}
	}

	value, latest := m.getValue(key)
	return protomsg.Clone(value), latest, nil
}

func (m *Map[T]) getValue(key string) (T, string) {
	old, ok := m.data[key]
	if ok {
		return old.value, old.version
	} else {
		var zero T
		return zero, missing
	}
}

// hasChanged returns true if the latest version of the key is different from version.
func (m *Map[T]) hasChanged(key string, version string) bool {
	_, latest := m.getValue(key)
	return version != latest
}
