package cuckooc

// actionMultiplexer is used to fetch the appropriate handler for a given action
var actionMultiplexer map[string]func(f *filterWrapper, args []string) (result string, err error)
