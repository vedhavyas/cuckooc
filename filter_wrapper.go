package cuckooc

import (
	"context"

	"github.com/vedhavyas/cuckoo-filter"
)

// filterWrapper is a wrapper over cuckoo filter
type filterWrapper struct {
	f     *cuckoo.Filter
	cmdCh <-chan Executor
}

// newWrapper returns a new filter wrapper.
// filter will not be initialised since the caller will be
// a filter manager passing cmds to appropriate filter wrapper which
// is running in its own go routine. We do not want to block manager to create
// filter. Hence, we off load it to filter wrapper's go routine
func newWrapper(cmdCh <-chan Executor) *filterWrapper {
	return &filterWrapper{cmdCh: cmdCh}
}

// listen will starts to listen commands over cmdCh channel
// blocking call. should be called in a different go-routine
func (fw *filterWrapper) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// TODO: backup and exit ?
			return
		case exe := <-fw.cmdCh:
			result, err := exe.Execute(fw)
			exe.Respond(result, err)
		}
	}
}
