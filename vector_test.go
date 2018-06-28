package goapl

import "testing"

func TestVectorAddDelete(t *testing.T) {
	vctInts := NewVector()
	for i:= 0; i < 10; i++ {
		vctInts.Push(uint64(i))
	}

	vctInts.RemoveIndex(0)
	n := vctInts.Length()
	if n != 9 {
		t.Fatal("ERR:n should be 9")
	}

	v, err := vctInts.Get(0)
	if err != nil || v != 1 {
		t.Fatal("ERR:index 0 should be 1")
	}
	
}