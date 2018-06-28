package goapl

import (
	"errors"
	"fmt"
)

type Vector struct {
	a   []uint64
	num int
}

func NewVector() *Vector {
	arr := make([]uint64, 0)
	return &Vector{a: arr, num:0}
}

func (v *Vector) Push(n uint64) {
	v.a = append(v.a, n)
	v.num ++
}

func (v *Vector) Get(index int) (uint64, error) {
	if index > len(v.a)-1 {
		return 0, errors.New("ERR:overflow")
	}

	return v.a[index], nil
}

func (v *Vector) RemoveIndex(index int) (uint64, error) {
	if index > len(v.a)-1 {
		return 0, errors.New("ERR:overflow")
	}

	n := v.a[index]
	v.a = append(v.a[:index], v.a[index+1:]...)
	v.num --
	return n, nil
}

func (v *Vector) RemoveValue(n uint64) {

	for i := 0; i < v.num; i++ {
		if v.a[i] != n {
			continue
		}
		v.a = append(v.a[:i], v.a[i+1:]...)
		v.num = v.num - 1
	}

	if len(v.a) != v.num {
		fmt.Println("ERR")
	}
}

func (v *Vector) Length() (n int) {
	n = v.num
	return n
}
