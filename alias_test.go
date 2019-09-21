package main

import "testing"

func TestExpandAlias(t *testing.T) {
	setAlias("ls", "ls -l")
	input, ans := "ls -h", "ls -l -h"

	if expandAlias(input) != ans {
		t.Fail()
	}

	setAlias("gst", "git status")
	input, ans = "gst", "git status"

	if expandAlias(input) != ans {
		t.Fail()
	}
}
