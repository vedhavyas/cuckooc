package cuckooc

import "github.com/vedhavyas/cuckoo-filter"

// filterWrapper is a wrapper over cuckoo filter
type filterWrapper struct {
	f      *cuckoo.Filter
	execCh chan Executor
}
