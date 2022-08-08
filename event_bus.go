package evbus

import (
	"fmt"
	"reflect"
	"sync"
)

type Bus interface {
	Subscribe(topic string, function interface{}) error
	UnSubscribe(topic string, function interface{}) error
	Publish(topic string, args ...interface{}) error
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

func New() Bus {
	return &EvBus{
		handlers: make(map[string][]*evHandler),
		lock:     sync.RWMutex{},
	}
}

func (b *EvBus) Subscribe(topic string, function interface{}) error {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(function))
	}

	reflectFunc := reflect.ValueOf(function)

	if h, _ := b.findHandler(topic, reflectFunc); h != nil {
		return fmt.Errorf("this function is already registered in %s", topic)
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	handler := &evHandler{
		stream:   make(chan []reflect.Value),
		callback: reflectFunc,
	}

	b.handlers[topic] = append(b.handlers[topic], handler)

	go func() {
		for v := range handler.stream {
			handler.callback.Call(v)
			b.wg.Done()
		}
	}()

	return nil
}

func (b *EvBus) Publish(topic string, args ...interface{}) error {
	hs, ok := b.handlers[topic]
	if !ok {
		return fmt.Errorf("%s is an unregistered topic", topic)
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	hs = append([]*evHandler{}, hs...)

	b.wg.Add(len(b.handlers[topic]))
	go func(handlers []*evHandler, args []interface{}) {
		for _, h := range handlers {
			h.stream <- b.convertHandlerArgs(args)
		}
	}(hs, args)

	return nil
}

func (b *EvBus) UnSubscribe(topic string, function interface{}) error {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(function))
	}

	h, i := b.findHandler(topic, reflect.ValueOf(function))
	if h == nil {
		return fmt.Errorf("this function is not registered in the %s", topic)
	}

	close(h.stream)
	b.deleteTargetHandler(i, topic)

	return nil
}

func (b *EvBus) Wait() {
	b.wg.Wait()
}

func (b *EvBus) deleteTargetHandler(index int, topic string) {
	b.handlers[topic] = append(b.handlers[topic][:index], b.handlers[topic][index+1:]...)
}

func (b *EvBus) convertHandlerArgs(args []interface{}) (hArgs []reflect.Value) {
	hArgs = make([]reflect.Value, len(args))

	for i, arg := range args {
		hArgs[i] = reflect.ValueOf(arg)
	}

	return hArgs
}

func (b *EvBus) findHandler(topic string, fn reflect.Value) (*evHandler, int) {
	if _, ok := b.handlers[topic]; !ok {
		return nil, -1
	}

	for i, h := range b.handlers[topic] {
		if h.callback.Pointer() == fn.Pointer() {
			return h, i
		}
	}

	return nil, -1
}
