package cuckooc

import (
	"reflect"
	"testing"

	"github.com/vedhavyas/cuckoo-filter"
)

func Test_createHandler(t *testing.T) {
	tests := []struct {
		args  []string
		count uint32
		bs    uint8
		err   bool
	}{
		{
			args:  []string{"100", "16"},
			count: 8,
			bs:    16,
		},

		{
			args:  []string{"100"},
			count: 16,
			bs:    8,
		},

		{
			args:  []string{},
			count: 524288,
			bs:    8,
		},

		{
			args: []string{"100", "18"},
			err:  true,
		},

		{
			args: []string{"not a number"},
			err:  true,
		},

		{
			args: []string{"100", "not a number"},
			err:  true,
		},
	}

	for _, c := range tests {
		fw := new(filterWrapper)
		_, err := createHandler(fw, c.args)
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
		args   []string
		result string
		err    bool
	}{
		{
			args:   []string{"x", "Y", "Z", "abc"},
			result: "true true true true",
		},

		{
			args:   []string{"a", "b"},
			result: "true true",
		},

		{
			args: []string{},
			err:  true,
		},
	}

	fw := &filterWrapper{f: cuckoo.StdFilter(), cmdCh: nil}
	for _, c := range tests {
		r, err := setHandler(fw, c.args)
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

func Test_setUniqueHandler(t *testing.T) {
	tests := []struct {
		args   []string
		result string
		err    bool
	}{
		{
			args:   []string{"x", "Y", "Z", "abc"},
			result: "true true true true",
		},

		{
			args:   []string{"a", "b", "x", "y"},
			result: "true true true true",
		},

		{
			args: []string{},
			err:  true,
		},
	}

	fw := &filterWrapper{f: cuckoo.StdFilter(), cmdCh: nil}
	for _, c := range tests {
		r, err := setUniqueHandler(fw, c.args)
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

func Test_checkHandler(t *testing.T) {
	setArgs := []string{"a", "1", "x", "Y", "X", "abc", "test"}
	tests := []struct {
		args   []string
		result string
		err    bool
	}{
		{
			args:   []string{"a", "b", "c"},
			result: "true false false",
		},

		{
			args:   []string{"1", "x", "Y", "ABC"},
			result: "true true true false",
		},

		{
			args: []string{},
			err:  true,
		},
	}

	fw := &filterWrapper{f: cuckoo.StdFilter(), cmdCh: nil}
	setUniqueHandler(fw, setArgs)
	for _, c := range tests {
		result, err := checkHandler(fw, c.args)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if result != c.result {
			t.Fatalf("expected %s but got %s", c.result, result)
		}
	}
}

func Test_deleteHandler(t *testing.T) {
	setArgs := []string{"a", "1", "x", "Y", "X", "abc", "test"}
	tests := []struct {
		args   []string
		result string
		err    bool
	}{
		{
			args:   []string{"a", "b", "c"},
			result: "true false false",
		},

		{
			args:   []string{"1", "x", "Y", "ABC"},
			result: "true true true false",
		},

		{
			args: []string{},
			err:  true,
		},
	}

	fw := &filterWrapper{f: cuckoo.StdFilter(), cmdCh: nil}
	setUniqueHandler(fw, setArgs)
	for _, c := range tests {
		result, err := deleteHandler(fw, c.args)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if result != c.result {
			t.Fatalf("expected %s but got %s", c.result, result)
		}
	}
}
