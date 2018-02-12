package cuckooc

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vedhavyas/cuckoo-filter"
)

const (
	ok    = "true"
	notOk = "false"
)

// handlerMux is used to fetch the appropriate handler for a given action
var handlerMux = map[string]func(config Config, f *filter, args []string) (result string, err error){
	"create":     createHandler,
	"set":        setHandler,
	"setu":       setUniqueHandler,
	"check":      checkHandler,
	"delete":     deleteHandler,
	"count":      countHandler,
	"loadfactor": loadFactorHandler,
	"backup":     backupHandler,
	"reload":     reloadHandler,
}

// createHandler creates cuckoo filter if not created already
// error when filter is already created
//
// args for create handler
// [filter-name] create [count] [bucket size]
// if count/bucket size are not provide, defaults to standard cuckoo filter
func createHandler(_ Config, f *filter, args []string) (result string, err error) {
	if f.f != nil {
		return "", fmt.Errorf("filter already exists")
	}

	var count uint32 = 4 << 20
	var bs uint8 = 8
	if len(args) >= 1 {
		c, err := strconv.Atoi(args[0])
		if err != nil {
			return result, fmt.Errorf("not a valid count: %v", args[0])
		}
		count = uint32(c)

		if len(args) == 2 {
			c, err := strconv.Atoi(args[1])
			if err != nil {
				return result, fmt.Errorf("not a valid bucket size: %v", args[1])
			}
			bs = uint8(c)
		}
	}

	filter, err := cuckoo.NewFilterWithBucketSize(count, bs)
	if err != nil {
		return result, fmt.Errorf("failed to create filter: %v", err)
	}

	f.f = filter
	return ok, nil
}

func commonHandler(f func([]byte) bool, args []string) (result string, err error) {
	if len(args) < 1 {
		return result, fmt.Errorf("require atleast one argument")
	}

	var results []string
	for _, x := range args {
		ok := f([]byte(x))
		results = append(results, fmt.Sprint(ok))
	}

	return strings.Join(results, " "), nil
}

// setHandler handles the set operations on the filter
//
// cmd format for setHandler
// [filter-name] set [args...]
// handler can handle multiple set operations in a single command
// requires at least one argument
func setHandler(_ Config, f *filter, args []string) (result string, err error) {
	return commonHandler(f.f.UInsert, args)
}

// setUniqueHandler handles the set unique operations
//
// format for setUniqueHandler
// [filter-name] setu [args...]
// requires at least one argument
func setUniqueHandler(_ Config, f *filter, args []string) (result string, err error) {
	return commonHandler(f.f.UInsertUnique, args)

}

// checkHandler handles the lookup operations
//
// format for checkHandler
// [filter-name] check [args...]
// requires at least one argument
func checkHandler(_ Config, f *filter, args []string) (result string, err error) {
	return commonHandler(f.f.ULookup, args)
}

// deleteHandler handles delete operations
//
// format for deleteHandler
// [filter-name] delete [args...]
// requires at least one argument
func deleteHandler(_ Config, f *filter, args []string) (result string, err error) {
	return commonHandler(f.f.UDelete, args)
}

// countHandler handles the count of items set in filter
//
// format for countHandler
// [filter-name] count
// any args passed will be ignored
func countHandler(_ Config, f *filter, _ []string) (result string, err error) {
	return fmt.Sprint(f.f.UCount()), nil
}

// loadFactorHandler handles requests for the load factor of a filter
//
// format for loadFactorHandler
// [filter-name] loadfactor
//any args passed will be ignored
func loadFactorHandler(_ Config, f *filter, _ []string) (result string, err error) {
	return fmt.Sprintf("%.4f", f.f.ULoadFactor()), nil
}

// backupHandler handles the backup requests for filters
//
// format for backupHandler
// [filter-name] backup [path to backup folder(overrides the one provided in config)]
func backupHandler(config Config, f *filter, args []string) (result string, err error) {
	path := config.BackupFolder
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		path = args[0]
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return result, fmt.Errorf("backup folder not provided")
	}

	// create the folder if not exists
	err = os.MkdirAll(path, 0700)
	if err != nil {
		return result, fmt.Errorf("failed to create backup directory: %v", err)
	}

	// let's encode the filter
	var buf bytes.Buffer
	err = f.f.Encode(&buf)
	if err != nil {
		return result, fmt.Errorf("failed to encode filter %s: %v", f.name, err)
	}

	bw := backupFilter{Name: f.name, FilterBytes: buf.Bytes()}
	path = filepath.Join(path, fmt.Sprintf("%s-latest.cf", f.name))
	fh, err := os.Create(path)
	if err != nil {
		return result, fmt.Errorf("failed to create backup file: %v", err)
	}
	defer fh.Close()

	enc := gob.NewEncoder(fh)
	err = enc.Encode(bw)
	if err != nil {
		return result, fmt.Errorf("failed to backup the filter: %v", err)
	}

	// ensure data is committed to the storage
	err = fh.Sync()
	if err != nil {
		return result, fmt.Errorf("failed to sync the file: %v", err)
	}

	return ok, nil
}

// reloadHandler handles the requests to load the filter from last backup
//
// format for reload handler
// [filter-name] reload [path to backup folder(overrides the one provided in config)]
func reloadHandler(config Config, f *filter, args []string) (result string, err error) {
	path := config.BackupFolder
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		path = args[0]
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return result, fmt.Errorf("backup folder not set to load filter from")
	}

	path = filepath.Join(path, fmt.Sprintf("%s-latest.cf", f.name))
	fh, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("failed to read backup: %v", err)
	}
	defer fh.Close()

	var bw backupFilter
	rd := bufio.NewReader(fh)
	dec := gob.NewDecoder(rd)
	err = dec.Decode(&bw)
	if err != nil {
		return result, fmt.Errorf("failed to decode filter: %v", err)
	}

	f.f, err = cuckoo.Decode(bytes.NewReader(bw.FilterBytes))
	if err != nil {
		return result, fmt.Errorf("failed to load filter: %v", err)
	}

	return ok, nil
}
