package observer

import (
	"slices"
	"sync"

	"golang.org/x/exp/maps"
)

// UnsubscribeFunc denotes a function used for unsubscribing from a channel.
type UnsubscribeFunc = func()

// ReadonlyChannel is an observable channel which can only be subscribed to without being able to publish events.
type ReadonlyChannel[T any] interface {
	// Subscribe subscribes the observer-function to published events of this channel. Once an event is
	// published, the observer-function will be called with the published value. To unsubscribe again, run the returned
	// UnsubscribeFunc.
	Subscribe(observer func(T)) UnsubscribeFunc
}

// ReadWriteChannel is an observable channel which can be used for publishing and subscribing.
type ReadWriteChannel[T any] interface {
	ReadonlyChannel[T]
	// Publish publishes a new event and notifies all subscribers.
	Publish(data T)
	// Clear force-unsubscribes all subscribers from the channel.
	Clear()
}

type channel[T any] struct {
	callbacks map[uint64]func(T)
	nextID    uint64
	isAsync   bool
	mu        *sync.RWMutex
}

// NewSyncChannel create a new ReadWriteChannel which notifies subscribers sequentially on event publishing. Long
// executing subscribers will block the execution of the remaining subscribers.
func NewSyncChannel[T any]() ReadWriteChannel[T] {
	return &channel[T]{
		callbacks: make(map[uint64]func(T)),
		isAsync:   false,
		mu:        &sync.RWMutex{},
	}
}

// NewAsyncChannel creates a new ReadWriteChannel which notifies each subscriber asynchronously. Every subscriber
// will be notified in a dedicated Go-routine thus preventing long-execution-blocking.
func NewAsyncChannel[T any]() ReadWriteChannel[T] {
	return &channel[T]{
		callbacks: make(map[uint64]func(T)),
		isAsync:   true,
		mu:        &sync.RWMutex{},
	}
}

func (o *channel[T]) Subscribe(observer func(T)) UnsubscribeFunc {
	o.mu.Lock()
	defer o.mu.Unlock()

	id := o.nextID
	o.nextID++
	o.callbacks[id] = observer
	return func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		delete(o.callbacks, id)
	}
}

func (o *channel[T]) Publish(data T) {
	o.mu.RLock()
	isAsync := o.isAsync
	callbacks := make(map[uint64]func(T), len(o.callbacks))
	maps.Copy(callbacks, o.callbacks)
	o.mu.RUnlock()

	callbackIDs := maps.Keys(callbacks)
	slices.Sort(callbackIDs)

	for _, callbackID := range callbackIDs {
		observer := o.callbacks[callbackID]
		if isAsync {
			go observer(data)
		} else {
			observer(data)
		}
	}
}

func (o *channel[T]) Clear() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.callbacks = make(map[uint64]func(T))
}
