package cuckooc

import (
	"context"
	"log"
	"os"
	"sync"
)

// gkLog with specific prefix set
var gkLog = log.New(os.Stderr, "GK: ", log.LstdFlags)

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
func newFilterWrapper(parent context.Context, name string, config Config, wg *sync.WaitGroup, gkCmd chan<- string) *filterWrapper {
	ctx, cancel := context.WithCancel(parent)
	cmdCh := make(chan Executor)
	filter := newFilter(name, cmdCh, gkCmd)
	wg.Add(1)
	go filter.listen(ctx, config, wg)
	return &filterWrapper{cancel: cancel, cmdCh: cmdCh}
}

// Gatekeeper is the switching point for all the filter requests
//
// workflow is that cmd is sent over to Gatekeeper, which will route it appropriate
// filter, creates one if not available.
type Gatekeeper struct {
	CMDCh   chan Executor
	gkCmd   chan string
	filters map[string]*filterWrapper
}

// NewGatekeeper for a new Gatekeeper
func NewGatekeeper(cmdCh chan Executor) *Gatekeeper {
	return &Gatekeeper{
		CMDCh:   cmdCh,
		filters: make(map[string]*filterWrapper),
		gkCmd:   make(chan string),
	}
}

// Start initiates gatekeeper to listen for commands to route over cmdCh
// blocking call, will need to start in a separate go routine
func (gk *Gatekeeper) Start(ctx context.Context, config Config, wg *sync.WaitGroup) {
	gkLog.Println("Starting GateKeeper...")
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			gkLog.Println("Stopping GateKeeper...")
			return
		case cmd := <-gk.CMDCh:
			gkLog.Printf("new request for filter %s\n", cmd.FilterName())
			fw, ok := gk.filters[cmd.FilterName()]
			if !ok {
				gkLog.Printf("creating a new filter wrapper: %s\n", cmd.FilterName())
				fw = newFilterWrapper(ctx, cmd.FilterName(), config, wg, gk.gkCmd)
				gk.filters[cmd.FilterName()] = fw
			}

			// lets not wait here for filter wrapper
			// send the command over a different go routine
			go func(cmdCh chan<- Executor, cmd Executor) {
				cmdCh <- cmd
				gkLog.Printf("request routed to filter: %s\n", cmd.FilterName())
			}(fw.cmdCh, cmd)
		case filterName := <-gk.gkCmd:
			gkLog.Printf("stop request from filter %s\n", filterName)
			fw := gk.filters[filterName]
			fw.cancel()
			delete(gk.filters, filterName)
		}
	}
}
