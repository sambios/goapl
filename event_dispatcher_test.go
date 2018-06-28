package goapl

import (
	"testing"
	"time"
)

const HELLO_WORLD = "helloWorld"

func TestDispatcher(t *testing.T) {
	var x int
	dispatcher := NewEventDispatcher()

	x = 1
	dispatcher.AddEventListener(1, func(){
		x++
	})

	time.Sleep(time.Second * 2)
	//dispatcher.RemoveEventListener(HELLO_WORLD, listener)

	dispatcher.EventTrigger(1)
	if x != 2 {
		t.Error("Event listener not exec!")
	}

}

