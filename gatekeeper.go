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

// Gatekeeper is the switching point for all the filter requests
//
// workflow is that cmd is sent over to Gatekeeper, which will route it appropriate
// filter, creates one if not available.
type Gatekeeper struct {
	cmdCh   <-chan Executor
	filters map[string]*filterWrapper
}

// NewGatekeeper for a new Gatekeeper
func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{
		cmdCh:   make(chan Executor),
		filters: make(map[string]*filterWrapper),
	}
}

// Start initiates gatekeeper to listen for commands to route over cmdCh
// blocking call, will need to start in a separate go routine
func (gk *Gatekeeper) Start(ctx context.Context, config Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-gk.cmdCh:
			// TODO check for keeper level actions
			fw, ok := gk.filters[cmd.FilterName()]
			if !ok {
				fw = newFilterWrapper(ctx, cmd.FilterName(), config, wg)
				gk.filters[cmd.FilterName()] = fw
			}

			// lets not wait for wrapper to get free
			// send the command over a different go routine
			go func(cmdCh chan<- Executor, cmd Executor) {
				cmdCh <- cmd
			}(fw.cmdCh, cmd)
		}
	}
}
