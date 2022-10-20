package main

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	err := run()
	fmt.Println(err)
	if err != nil {
		t.Error("failed run()")
	}
}
