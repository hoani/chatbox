package hal

import (
	"fmt"
	"sync"
)

type button struct {
	value bool
	lock  sync.Mutex
}

func (b *button) start() {
	go func() {
		for {
			fmt.Scanln()
			b.lock.Lock()
			b.value = !b.value
			b.lock.Unlock()
		}
	}()
}

func (b *button) get() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.value
}
