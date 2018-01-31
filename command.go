package cuckooc

import (
	"fmt"
	"strings"
)

// Executor executes command/command-set on a given filter
// and responds with a result
type Executor interface {
	FilterName() string
	Execute(f *filterWrapper) (result string, err error)
	Respond(result string, err error)
}

// command refers to single command to be performed on the filter using args provided
// and a Response channel to receive the execution result
type command struct {
	Filter       string
	Action       string
	Args         []string
	ResponseChan chan<- string
}

// FilterName returns the name of the filter
func (i command) FilterName() string {
	return i.Filter
}

// Respond sends the result/error over the response chan
// false if error
func (i command) Respond(result string, err error) {
	if err != nil {
		result = fmt.Sprintf("false(%v)", err)
	}

	i.ResponseChan <- result
}

// Execute fetches the appropriate action handler and executes the action on filter
func (i command) Execute(f *filterWrapper) (result string, err error) {
	ah, ok := actionMultiplexer[i.Filter]
	if !ok {
		return result, fmt.Errorf("unknown action: %s", i.Action)
	}

	return ah(f, i.Args)
}

// parseCommand parses the cmd and returns an command
//
// Format of the command is as follows
// [filter name] [action] [args...]
// command requires at least filter-name and an action
func parseCommand(cmd string, respChan chan<- string) (*command, error) {

	var args []string
	for _, arg := range strings.Split(strings.TrimSpace(cmd), " ") {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		args = append(args, arg)
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("require atleast 2 arguments")
	}

	return &command{
		Filter:       strings.ToLower(args[0]),
		Action:       strings.ToLower(args[1]),
		Args:         args[2:],
		ResponseChan: respChan,
	}, nil
}
