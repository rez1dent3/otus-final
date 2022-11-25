package bus

type CommandBusInterface interface {
	Fire(string, any)
	Subscribe(string, func(any))
}

type impl struct {
	cmd map[string][]func(any)
}

func NewSyncBus() CommandBusInterface {
	return &impl{cmd: make(map[string][]func(any))}
}

func (c *impl) Fire(name string, value any) {
	if fnList, ok := c.cmd[name]; ok {
		for _, fn := range fnList {
			fn(value)
		}
	}
}

func (c *impl) Subscribe(name string, fn func(any)) {
	if _, ok := c.cmd[name]; !ok {
		c.cmd[name] = make([]func(any), 1)
		c.cmd[name][0] = fn
		return
	}

	c.cmd[name] = append(c.cmd[name], fn)
}
