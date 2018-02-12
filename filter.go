package cuckooc

import (
	"context"

	"github.com/vedhavyas/cuckoo-filter"
)

// filter is a wrapper over cuckoo filter
type filter struct {
	name  string
	f     *cuckoo.Filter
	cmdCh <-chan Executor
}

// backupFilter is used only to backup the filter contents to a persistent disk
// Since cuckoo filter provides the option to encode its contents through gob encoder
// We use those bytes along with other meta data required for the filter wrapper
type backupFilter struct {
	Name        string
	FilterBytes []byte
}

// newWrapper returns a new filter wrapper.
// filter will not be initialised since the caller will be
// a filter manager passing commands to appropriate filter wrapper which
// is running in its own go routine. We do not want to block manager to create
// filter. Hence, we off load it to filter wrapper's go routine
func newWrapper(name string, cmdCh <-chan Executor) *filter {
	return &filter{name: name, cmdCh: cmdCh}
}

// listen will starts to listen commands over cmdCh channel
// blocking call. should be called in a different go-routine
func (f *filter) listen(ctx context.Context, config Config) {
	for {
		select {
		case <-ctx.Done():
			// TODO: backup and exit ?
			return
		case exe := <-f.cmdCh:
			result, err := exe.Execute(config, f)
			exe.Respond(result, config.Debug, err)
		}
	}
}
