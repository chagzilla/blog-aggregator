package main

type command struct {
	name string
	arguments []string
}

type commands struct {
	callbacks map[string]func(*state, command) error
}

func (c * commands) run(s * state, com command) error {
	return c.callbacks[com.name](s, com)
}

func (c * commands) register(name string, f func(*state, command) error) {
	c.callbacks[name] = f
}

