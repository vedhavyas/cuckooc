package cuckooc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vedhavyas/cuckoo-filter"
)

// actionMultiplexer is used to fetch the appropriate handler for a given action
var actionMultiplexer = map[string]func(fw *filterWrapper, args []string) (result string, err error){
	"create": createHandler,
	"set":    setHandler,
	"setu":   setUniqueHandler,
}

// createHandler creates cuckoo filter if not created already
// error when filter is already created
//
// args for create handler
// [filter-name] create [count] [bucket size]
// if count/bucket size are not provide, defaults to standard cuckoo filter
func createHandler(fw *filterWrapper, args []string) (result string, err error) {
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

// setHandler handles the set operations on the filter
//
// cmd format for setHandler
// [filter-name] set [args...]
// handler can handle multiple set operations in a single command
// requires at least one argument
func setHandler(fw *filterWrapper, args []string) (result string, err error) {
	if len(args) < 1 {
		return result, fmt.Errorf("require atleast one arg for set operation")
	}

	var results []string
	for _, x := range args {
		ok := fw.f.UInsert([]byte(x))
		results = append(results, fmt.Sprint(ok))
	}

	return strings.Join(results, " "), nil
}

// setUniqueHandler handles the set unique operations
//
//format for setUniqueHandler
// [filter-name] setu [args...]
// requires at least one argument
func setUniqueHandler(fw *filterWrapper, args []string) (result string, err error) {
	if len(args) < 1 {
		return result, fmt.Errorf("require atleast one arg")
	}

	var results []string
	for _, x := range args {
		ok := fw.f.UInsertUnique([]byte(x))
		results = append(results, fmt.Sprint(ok))
	}

	return strings.Join(results, " "), nil

}
