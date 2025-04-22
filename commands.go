package main

type commands struct {
	dict map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.dict[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	err := c.dict[cmd.name](s, cmd)

	return err
}
