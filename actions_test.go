package cuckooc

import (
	"reflect"
	"testing"

	"github.com/vedhavyas/cuckoo-filter"
)

func Test_createHandler(t *testing.T) {
	tests := []struct {
		cmd   string
		count uint32
		bs    uint8
		err   bool
	}{
		{
			cmd:   "test create 100 16",
			count: 8,
			bs:    16,
		},

		{
			cmd:   "test create 100",
			count: 16,
			bs:    8,
		},

		{
			cmd:   "test create",
			count: 524288,
			bs:    8,
		},

		{
			cmd: "test create 100 18",
			err: true,
		},

		{
			cmd: "test create abs",
			err: true,
		},

		{
			cmd: "test create 100 abs",
			err: true,
		},
	}

	for _, c := range tests {
		i, err := parseCommand(c.cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		fw := new(filterWrapper)
		_, err = createHandler(fw, i.Args)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		fr := reflect.ValueOf(fw.f).Elem()
		count := fr.FieldByName("totalBuckets")
		if c.count != uint32(count.Uint()) {
			t.Fatalf("expected %d count but got %d", c.count, count.Uint())
		}

		bs := fr.FieldByName("bucketSize")
		if c.bs != uint8(bs.Uint()) {
			t.Fatalf("expected %d bucket size but got %d", c.bs, bs.Uint())
		}
	}
}

func Test_setHandler(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
		err    bool
	}{
		{
			cmd:    "test set x Y Z abc",
			result: "true true true true",
		},

		{
			cmd:    "test set a  b",
			result: "true true",
		},

		{
			cmd: "test set",
			err: true,
		},
	}

	fw := &filterWrapper{f: cuckoo.StdFilter(), cmdCh: nil}
	for _, c := range tests {
		i, err := parseCommand(c.cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		r, err := setHandler(fw, i.Args)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if r != c.result {
			t.Fatalf("expected %s but got %s", c.result, r)
		}
	}
}
