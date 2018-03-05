package cuckooc

import (
	"reflect"
	"testing"
)

func Test_newInstruction(t *testing.T) {
	tests := []struct {
		cmd    string
		filter string
		action string
		args   []string
		err    bool
	}{
		{
			cmd:    "test create 1 2 3 4 ",
			filter: "test",
			action: "create",
			args:   []string{"1", "2", "3", "4"},
		},

		{
			cmd:    "test set x Y z",
			filter: "test",
			action: "set",
			args:   []string{"x", "Y", "z"},
		},

		{
			cmd:    "test set x Y  z",
			filter: "test",
			action: "set",
			args:   []string{"x", "Y", "z"},
		},

		{
			cmd:    "test set x     Y  z",
			filter: "test",
			action: "set",
			args:   []string{"x", "Y", "z"},
		},

		{
			cmd: "test ",
			err: true,
		},
	}

	for _, c := range tests {
		i, err := parseCommand(c.cmd, nil)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if i.Filter != c.filter {
			t.Fatalf("expected %s filter but got %s", c.filter, i.Filter)
		}

		if i.Action != c.action {
			t.Fatalf("expected %s action but got %s", c.action, i.Action)
		}

		if !reflect.DeepEqual(i.Args, c.args) {
			t.Fatalf("expected %v args but got %v", c.args, i.Args)
		}
	}
}

func TestCommand_readCommands(t *testing.T) {
	tests := []struct {
		s  string
		ss []string
	}{
		{
			ss: nil,
		},
		{
			s:  "test new",
			ss: []string{"test new"},
		},

		{
			s:  "test new\n test1 set x a b c",
			ss: []string{"test new", " test1 set x a b c"},
		},

		{
			s:  "test new\ntest1 set x a b c",
			ss: []string{"test new", "test1 set x a b c"},
		},
	}

	for _, c := range tests {
		ss := readCommands([]byte(c.s))
		if !reflect.DeepEqual(c.ss, ss) {
			t.Fatalf("expected %v but got %v", c.ss, ss)
		}
	}
}
