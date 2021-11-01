package main

import (
	"os"
	"testing"
)

func TestKeepOnArg(t *testing.T) {
	cases := []struct {
		args []string
		want bool
	}{
		{[]string{"r~"}, false},
		{[]string{"r~", "--keep-going"}, true},
		{[]string{"r~", "-keep-going"}, true},
		{[]string{"r~", "-k"}, true},
		{[]string{"r~", "--k"}, true}, // We'll allow this one...
		{[]string{"-keep-going", "this", "is", "silly"}, false},
		{[]string{"r~", "-i", "-r", "-y", "3", "-k"}, true},
		{[]string{"r~", "-i", "-r", "-y", "3"}, false},
		{[]string{"r~", "find", "me", "-keep-going", "in", "the", "middle"}, true},
		{[]string{"r~", "-k", "-r", "--keep-going"}, true},
	}
	for _, c := range cases {
		os.Args = c.args
		if keepOnArg() != c.want {
			t.Errorf("args %v, wanted %t\n", c.args, c.want)
		}
	}
}
