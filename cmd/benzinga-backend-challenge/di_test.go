package main

import "testing"

func TestDI(t *testing.T) {
	dic := newDIContainer(&flags{})
	_, err := dic.httpHandler()
	if err != nil {
		t.Fatal(err)
	}
}
