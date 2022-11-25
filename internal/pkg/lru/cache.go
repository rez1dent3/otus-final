package lru

import (
	"container/list"
	"sync"

	"github.com/rez1dent3/otus-final/internal/pkg/bus"
)

const EventEvict = "event_evict"

type CacheInterface interface {
	Put(string, any) bool
	Get(string) (any, bool)
	Has(string) bool
	Size() uint64
	Purge()
}

type entry struct {
	key string
	val any
}

type impl struct {
	size, limit uint64

	mu sync.RWMutex

	evict *list.List
	items map[string]*list.Element

	busCommand bus.CommandBusInterface
}

func New(sizeLimit uint64, busCommand bus.CommandBusInterface) CacheInterface {
	return &impl{limit: sizeLimit, evict: list.New(), busCommand: busCommand, items: make(map[string]*list.Element)}
}

func (c *impl) Put(key string, value any) bool {
	elementSize := c.elementSize(value)
	if c.limit < elementSize {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	defer c.removeOldest()

	if val, ok := c.items[key]; ok {
		c.evict.MoveToFront(val)
		ent := val.Value.(*entry)
		c.size -= c.elementSize(ent.val) - elementSize
		ent.val = value

		return true
	}

	ent := &entry{key: key, val: value}
	item := c.evict.PushFront(ent)

	c.items[key] = item
	c.size += elementSize

	return true
}

func (c *impl) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ent, ok := c.items[key]; ok {
		c.evict.MoveToFront(ent)

		return ent.Value.(*entry).val, true
	}

	return nil, false
}

func (c *impl) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.items[key]

	return ok
}

func (c *impl) Size() uint64 {
	return c.size
}

func (c *impl) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for c.evict.Back() != nil {
		c.delete(c.evict.Back())
	}
}

func (c *impl) delete(el *list.Element) {
	if el == nil {
		return
	}

	ent := el.Value.(*entry)

	c.size -= c.elementSize(ent.val)
	delete(c.items, ent.key)

	c.evict.Remove(el)

	c.busCommand.Fire(EventEvict, ent.val)
}

func (c *impl) removeOldest() {
	for c.limit < c.size && c.evict.Back() != nil {
		c.delete(c.evict.Back())
	}
}

// elementSize If the object implements the "size" method, then we calc the volume by this arg.
// Otherwise, we calculate the volume by the number of elements.
func (c *impl) elementSize(e any) uint64 {
	if val, ok := e.(interface {
		Size() uint64
	}); ok {
		return val.Size()
	}

	return 1
}
