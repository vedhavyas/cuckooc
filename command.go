package cuckooc

import (
	"fmt"
	"strings"
)

// Executor executes command/command-set on a given filter
// and responds with a result
type Executor interface {
	FilterName() string
	Execute(config Config, f *filter) (result string, err error)
	Respond(result string, debug bool, err error)
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
func (c command) FilterName() string {
	return c.Filter
}

// Respond sends the result/error over the response chan
// false if error
func (c command) Respond(result string, debug bool, err error) {
	if err != nil {
		result = notOk
		if debug {
			result = fmt.Sprintf("%s(%v)", notOk, err)
		}
	}

	c.ResponseChan <- result
}

// Execute fetches the appropriate action handler and executes the action on filter
func (c command) Execute(config Config, f *filter) (result string, err error) {
	ah, ok := handlerMux[c.Action]
	if !ok {
		return result, fmt.Errorf("unknown action: %s", c.Action)
	}

	return ah(config, f, c.Args)
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
