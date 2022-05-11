package main

import "testing"

func TestHelloWorld(t *testing.T) {
	t.Run("hello", func(t *testing.T) {
		t.Log("hello")
	})
}
