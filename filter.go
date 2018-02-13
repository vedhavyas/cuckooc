package cuckooc

import (
	"context"
	"fmt"
	"sync"

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

// newFilter returns a new filter wrapper.
// filter will not be initialised since the caller will be
// a filter manager passing commands to appropriate filter wrapper which
// is running in its own go routine. We do not want to block manager to create
// filter. Hence, we off load it to filter wrapper's go routine
func newFilter(name string, cmdCh <-chan Executor) *filter {
	return &filter{name: name, cmdCh: cmdCh}
}

// loadFilter will reload the last saved filter from persistent storage
// if load fails and if action is new, return nil
// else respond to command with error
func loadFilter(config Config, f *filter, action string) error {
	_, err := reloadHandler(config, f, nil)
	if err == nil {
		return nil
	}

	// load from the file system failed. check if the action is new, then create proceed
	// else fail
	if action == actionNew {
		return nil
	}

	return fmt.Errorf("filter doesn't exists")
}

// listen will starts to listen commands over cmdCh channel
// blocking call. should be called in a different go-routine
func (f *filter) listen(ctx context.Context, config Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// TODO: backup and exit ?
			// ensure f.f is not nil
			return
		case exe := <-f.cmdCh:
			if f.f == nil {
				err := loadFilter(config, f, exe.GetAction())
				if err != nil {
					exe.Respond("", config.Debug, err)
					continue
				}
			}

			result, err := exe.Execute(config, f)
			exe.Respond(result, config.Debug, err)

		}
	}
}
