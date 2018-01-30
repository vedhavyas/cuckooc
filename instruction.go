package cuckooc

import (
	"fmt"
	"strings"
)

// Executor executes instruction/instruction set on a given filter
// and responds with a result
type Executor interface {
	FilterName() string
	Execute(f *filterWrapper) (result string, err error)
	Respond(result string, err error)
}

// instruction refers to single instruction to be performed on the filter using args provided
// and a Response channel to receive the execution result
type instruction struct {
	Filter       string
	Action       string
	Args         []string
	ResponseChan chan<- string
}

// FilterName returns the name of the filter
func (i instruction) FilterName() string {
	return i.Filter
}

// Respond sends the result/error over the response chan
func (i instruction) Respond(result string, err error) {
	if err != nil {
		result = err.Error()
	}

	i.ResponseChan <- result
}

// Execute fetches the appropriate action handler and executes the action on filter
func (i instruction) Execute(f *filterWrapper) (result string, err error) {
	ah, ok := actionMultiplexer[i.Filter]
	if !ok {
		return result, fmt.Errorf("unknown action: %s", i.Action)
	}

	return ah(f, i.Args)
}

// newInstruction parses the cmd and returns an executor
//
// Format of the is as follows
// [filter name] [action] [args...]
func newInstruction(cmd string, respChan chan<- string) (*instruction, error) {
	pcmd := strings.Split(strings.TrimSpace(cmd), " ")
	if len(pcmd) < 2 {
		return nil, fmt.Errorf("require atleast 2 arguments")
	}

	return &instruction{
		Filter:       strings.ToLower(pcmd[0]),
		Action:       strings.ToLower(pcmd[1]),
		Args:         pcmd[2:],
		ResponseChan: respChan,
	}, nil
}
