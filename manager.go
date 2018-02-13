package cuckooc

import (
	"context"
	"sync"
)

// filterWrapper holds the bare minimum details for a filter
//
// cancel is a CancelFunc which can be used to stop the filter
// cmdCh to send in filter specific commands
type filterWrapper struct {
	cancel context.CancelFunc
	cmdCh  chan<- Executor
}

// newFilterWrapper creates a filter and starts it over a different go routine
//
// parent context is used to create a child context with parent
// name of the filter
// config of the current running service
// wg is the wait group for all go routines
func newFilterWrapper(parent context.Context, name string, config Config, wg *sync.WaitGroup) *filterWrapper {
	ctx, cancel := context.WithCancel(parent)
	cmdCh := make(chan Executor)
	filter := newFilter(name, cmdCh)
	wg.Add(1)
	go filter.listen(ctx, config, wg)
	return &filterWrapper{cancel: cancel, cmdCh: cmdCh}
}
