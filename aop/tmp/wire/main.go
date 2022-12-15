package main

import (
	"errors"
	"fmt"
	"time"
)

// Message is what greeters will use to greet guests.
type Message string

// NewMessage creates a default Message.
func NewMessage(phrase string) Message {
	return Message(phrase)
}

// NewGreeter initializes a Greeter. If the current epoch time is an even
// number, NewGreeter will create a grumpy Greeter.
func NewGreeter(m Message) Greeter {
	var grumpy bool
	if time.Now().Unix()%2 == 0 {
		grumpy = true
	}
	return Greeter{Message: m, Grumpy: grumpy}
}

// Greeter is the type charged with greeting guests.
type Greeter struct {
	Grumpy  bool
	Message Message
}

// Greet produces a greeting for guests.
func (g Greeter) Greet() Message {
	if g.Grumpy {
		return Message("Go away!")
	}
	return g.Message
}

// NewEvent creates an event with the specified greeter.
func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy {
		return Event{}, errors.New("could not create event: event greeter is grumpy")
	}
	return Event{Greeter: g}, nil
}

// Event is a gathering with greeters.
type Event struct {
	Greeter Greeter
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}
