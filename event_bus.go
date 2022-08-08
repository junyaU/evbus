package event_bus

import (
	"fmt"
	"reflect"
	"sync"
)

type EBus interface {
	Subscribe(topic string, function interface{}) error
	Publish(topic string, data ...interface{}) error
	Wait()
}

type EvBus struct {
	handlers map[string][]*evHandler
	lock     sync.RWMutex
	wg       sync.WaitGroup
}

type evHandler struct {
	stream   chan []reflect.Value
	callback reflect.Value
}

func New() EBus {
	return &EvBus{
		handlers: make(map[string][]*evHandler),
		lock:     sync.RWMutex{},
	}
}

func (b *EvBus) Subscribe(topic string, function interface{}) error {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(function))
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	h := &evHandler{
		stream:   make(chan []reflect.Value),
		callback: reflect.ValueOf(function),
	}

	b.handlers[topic] = append(b.handlers[topic], h)

	go func() {
		for v := range h.stream {
			h.callback.Call(v)
			b.wg.Done()
		}
	}()

	return nil
}

func (b *EvBus) Publish(topic string, data ...interface{}) error {
	hs, ok := b.handlers[topic]
	if !ok {
		return fmt.Errorf("%s is an unregistered topic", topic)
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	b.wg.Add(1)
	go func(handlers []*evHandler, args ...interface{}) {
		for _, h := range handlers {
			h.stream <- b.convertHandlerArgs(data)
		}
	}(hs, data)

	return nil
}

func (b *EvBus) UnSubscribe() error {

}

func (b *EvBus) Wait() {
	b.wg.Wait()
}

func (b *EvBus) convertHandlerArgs(args []interface{}) (hArgs []reflect.Value) {
	hArgs = make([]reflect.Value, len(args))

	for i, arg := range args {
		hArgs[i] = reflect.ValueOf(arg)
	}

	return hArgs
}
