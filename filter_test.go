package cuckooc

import (
	"context"
	"sync"
	"testing"
)

func runTests(t *testing.T, tests []struct {
	cmd    string
	result string
}, config Config) {

	cmdCh := make(chan Executor)
	respCh := make(chan string)
	f := newFilter("test", cmdCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go f.listen(ctx, config, wg)
	for _, c := range tests {
		exe, err := parseCommand(c.cmd, respCh)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmdCh <- exe
		result := <-respCh

		if result != c.result {
			t.Fatalf("expected result %s but got %s", c.result, result)
		}
	}
}

func Test_filterWrapper_integration(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test new 30",
			result: "true",
		},

		{
			cmd:    "test set a b c d e",
			result: "true true true true true",
		},

		{
			cmd:    "test setu a b f g",
			result: "true true true true",
		},

		{
			cmd:    "test check a b 1 2 c d e f g",
			result: "true true false false true true true true true",
		},

		{
			cmd:    "test count",
			result: "7",
		},

		{
			cmd:    "test loadfactor",
			result: "0.2188",
		},

		{
			cmd:    "test delete a c d f z",
			result: "true true true true false",
		},

		{
			cmd:    "test count",
			result: "3",
		},

		{
			cmd:    "test loadfactor",
			result: "0.0938",
		},
	}

	runTests(t, tests, Config{})
}

func TestFilter_loadFromFS_success(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test setu a b c d e f g",
			result: "true true true true true true true",
		},

		{
			cmd:    "test check a b 1 2 c d e f g",
			result: "true true false false true true true true true",
		},

		{
			cmd:    "test count",
			result: "7",
		},
	}

	runTests(t, tests, Config{BackupFolder: "./testdata/backups-2"})
}

func TestFilter_loadFromFS_failure(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test setu a b c",
			result: "false",
		},

		{
			cmd:    "test check a b",
			result: "false",
		},

		{
			cmd:    "test count",
			result: "false",
		},
	}

	runTests(t, tests, Config{BackupFolder: "./testdata/backups-3"})
}

func TestFilter_loadFromFS_action_success(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test new 100",
			result: "true",
		},

		{
			cmd:    "test setu a b c",
			result: "true true true",
		},

		{
			cmd:    "test check a b e",
			result: "true true false",
		},

		{
			cmd:    "test count",
			result: "3",
		},
	}

	runTests(t, tests, Config{BackupFolder: "./testdata/backups-3"})
}
