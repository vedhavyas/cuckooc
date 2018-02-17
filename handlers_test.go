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

	config := Config{}
	for _, c := range tests {
		f := new(filter)
		_, err := newHandler(config, f, c.args)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		fr := reflect.ValueOf(f.f).Elem()
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

	config := Config{}
	f := &filter{f: cuckoo.StdFilter(), cmdCh: nil}
	for _, c := range tests {
		r, err := setHandler(config, f, c.args)
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

	config := Config{}
	f := &filter{f: cuckoo.StdFilter(), cmdCh: nil}
	for _, c := range tests {
		r, err := setUniqueHandler(config, f, c.args)
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

	config := Config{}
	f := &filter{f: cuckoo.StdFilter(), cmdCh: nil}
	setUniqueHandler(config, f, setArgs)
	for _, c := range tests {
		result, err := checkHandler(config, f, c.args)
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

	config := Config{}
	f := &filter{f: cuckoo.StdFilter(), cmdCh: nil}
	setUniqueHandler(config, f, setArgs)
	for _, c := range tests {
		result, err := deleteHandler(config, f, c.args)
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

func Test_backupHandler(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		asConfig bool
		err      bool
	}{
		{
			name:     "test",
			path:     "./testdata/backups",
			asConfig: true,
			err:      false,
		},

		{
			name:     "test-1",
			path:     "./testdata/backups-1",
			asConfig: false,
			err:      false,
		},

		{
			name:     "test-2",
			path:     "",
			asConfig: true,
			err:      true,
		},
	}

	config := Config{}
	for _, c := range tests {
		var path string
		if c.asConfig {
			config.BackupFolder = c.path
			path = ""
		} else {
			path = c.path
			config.BackupFolder = ""
		}

		f := &filter{name: c.name, f: cuckoo.StdFilter()}
		_, err := setHandler(config, f, []string{"a", "b", "c", "d"})
		if err != nil {
			t.Fatalf("unexpected error for setHandler: %v", err)
		}

		res, err := backupHandler(config, f, []string{path})
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error for backupHandler: %v", err)
		}

		if res != "true" {
			t.Fatalf("expected true but got %s", res)
		}
	}
}

func Test_loadHandler(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		asConfig bool
		err      bool
	}{
		{
			name:     "test",
			path:     "./testdata/backups",
			asConfig: true,
			err:      false,
		},

		{
			name:     "test-1",
			path:     "./testdata/backups-1",
			asConfig: false,
			err:      false,
		},

		{
			name:     "test-2",
			path:     "",
			asConfig: true,
			err:      true,
		},
	}

	config := Config{}
	for _, c := range tests {
		var path string
		if c.asConfig {
			config.BackupFolder = c.path
			path = ""
		} else {
			path = c.path
			config.BackupFolder = ""
		}

		f := &filter{name: c.name}
		res, err := reloadHandler(config, f, []string{path})
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if res != "true" {
			t.Fatalf("expected true but got %s", res)
		}

		res, err = checkHandler(config, f, []string{"a", "b", "c", "d", "e"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res != "true true true true false" {
			t.Fatalf("expected \"true true true true false\" but got %s", res)
		}
	}
}

func Test_stopHandler(t *testing.T) {
	ch := make(chan string)
	f := newFilter("test", nil, ch)
	go func() {
		stopHandler(Config{}, f, nil)
	}()

	if f := <-ch; f != "test" {
		t.Fatalf("expected filter name \"test\" but got %s", f)
	}
}
