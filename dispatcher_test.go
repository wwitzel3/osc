package osc

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

// Test a successful method invocation.
func TestDispatcherDispatchOK(t *testing.T) {
	c := make(chan struct{})
	d := PatternMatching{
		"/bar": Method(func(msg Message) error {
			close(c)
			return nil
		}),
	}
	later := time.Now().Add(20 * time.Millisecond)
	b := Bundle{
		Timetag: FromTime(later),
		Packets: []Packet{
			Message{Address: "/bar"},
		},
	}
	if err := d.Dispatch(b, false); err != nil {
		t.Fatal(err)
	}
	<-c
}

// Test a method that returns an error.
func TestDispatcherDispatchError(t *testing.T) {
	d := PatternMatching{
		"/foo": Method(func(msg Message) error {
			return errors.New("oops")
		}),
	}
	later := time.Now().Add(20 * time.Millisecond)
	b := Bundle{
		Timetag: FromTime(later),
		Packets: []Packet{
			Message{Address: "/foo"},
		},
	}
	if err := d.Dispatch(b, false); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDispatcherDispatchNestedBundle(t *testing.T) {
	c := make(chan struct{})
	d := PatternMatching{
		"/foo": Method(func(msg Message) error {
			close(c)
			return nil
		}),
	}
	later := time.Now().Add(20 * time.Millisecond)
	b := Bundle{
		Timetag: FromTime(later),
		Packets: []Packet{
			Bundle{
				Timetag: FromTime(later.Add(20 * time.Millisecond)),
				Packets: []Packet{
					Message{Address: "/foo"},
				},
			},
		},
	}
	if err := d.Dispatch(b, false); err != nil {
		t.Fatal(err)
	}
	<-c
}

func TestDispatcherMiss(t *testing.T) {
	d := PatternMatching{
		"/foo": Method(func(msg Message) error {
			return nil
		}),
	}
	b := Bundle{
		Timetag: FromTime(time.Now()),
	}
	if err := d.Dispatch(b, false); err != nil {
		t.Fatal(err)
	}
}

func TestDispatcherInvoke(t *testing.T) {
	d := PatternMatching{
		"/foo": Method(func(msg Message) error {
			return errors.New("foo error")
		}),
		"/bar": Method(func(msg Message) error {
			return nil
		}),
	}
	msg := Message{Address: "/foo"}
	if err := d.Invoke(msg, false); err == nil {
		t.Fatal("expected error, got nil")
	}
	badMsg := Message{Address: "/["}
	if err := d.Invoke(badMsg, false); err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := d.Invoke(Message{Address: "/bar"}, false); err != nil {
		t.Fatal(err)
	}
	if err := d.Invoke(Message{Address: "/baz"}, false); err != nil {
		t.Fatal(err)
	}
	if err := d.invoke(badPacket{}, false); err == nil {
		t.Fatal("expected error, got nil")
	}
}
