package cuckooc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vedhavyas/cuckoo-filter"
)

// handlerMux is used to fetch the appropriate handler for a given action
// TODO(ved) backup handler
var handlerMux = map[string]func(config Config, fw *filterWrapper, args []string) (result string, err error){
	"create":     createHandler,
	"set":        setHandler,
	"setu":       setUniqueHandler,
	"check":      checkHandler,
	"delete":     deleteHandler,
	"count":      countHandler,
	"loadfactor": loadFactorHandler,
	"backup":     backupHandler,
}

// createHandler creates cuckoo filter if not created already
// error when filter is already created
//
// args for create handler
// [filter-name] create [count] [bucket size]
// if count/bucket size are not provide, defaults to standard cuckoo filter
func createHandler(_ Config, fw *filterWrapper, args []string) (result string, err error) {
	if fw.f != nil {
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

	f, err := cuckoo.NewFilterWithBucketSize(count, bs)
	if err != nil {
		return result, fmt.Errorf("failed to create filter: %v", err)
	}

	fw.f = f
	return "true", nil
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
func setHandler(_ Config, fw *filterWrapper, args []string) (result string, err error) {
	return commonHandler(fw.f.UInsert, args)
}

// setUniqueHandler handles the set unique operations
//
// format for setUniqueHandler
// [filter-name] setu [args...]
// requires at least one argument
func setUniqueHandler(_ Config, fw *filterWrapper, args []string) (result string, err error) {
	return commonHandler(fw.f.UInsertUnique, args)

}

// checkHandler handles the lookup operations
//
// format for checkHandler
// [filter-name] check [args...]
// requires at least one argument
func checkHandler(_ Config, fw *filterWrapper, args []string) (result string, err error) {
	return commonHandler(fw.f.ULookup, args)
}

// deleteHandler handles delete operations
//
// format for deleteHandler
// [filter-name] delete [args...]
// requires at least one argument
func deleteHandler(_ Config, fw *filterWrapper, args []string) (result string, err error) {
	return commonHandler(fw.f.UDelete, args)
}

// countHandler handles the count of items set in filter
//
// format for countHandler
// [filter-name] count
// any args passed will be ignored
func countHandler(_ Config, fw *filterWrapper, _ []string) (result string, err error) {
	return fmt.Sprint(fw.f.UCount()), nil
}

// loadFactorHandler handles requests for the load factor of a filter
//
// format for loadFactorHandler
// [filter-name] loadfactor
//any args passed will be ignored
func loadFactorHandler(_ Config, fw *filterWrapper, _ []string) (result string, err error) {
	return fmt.Sprintf("%.4f", fw.f.ULoadFactor()), nil
}

// backupHandler handles the backup requests for filters
//
// format for backupHandler
// [filter-name] backup [path to backup folder(overrides the one provided in config)]
func backupHandler(config Config, fw *filterWrapper, args []string) (result string, err error) {
	path := config.BackupFolder
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		path = args[0]
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return result, fmt.Errorf("backup folder not provided")
	}

	// create the folder if not exists
	err = os.MkdirAll(path, 0766)
	if err != nil {
		return result, fmt.Errorf("failed to create backup directory: %v", err)
	}

	// let's encode the filter
	var buf bytes.Buffer
	err = fw.f.Encode(&buf)
	if err != nil {
		return result, fmt.Errorf("failed to encode filter %s: %v", fw.name, err)
	}

	bw := backupFilter{Name: fw.name, FilterBytes: buf.Bytes()}
	path = filepath.Join(path, fmt.Sprintf("%s-latest.cf", fw.name))
	fh, err := os.Create(path)
	if err != nil {
		return result, fmt.Errorf("failed to create backu file: %v", err)
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

	return "true", nil
}
