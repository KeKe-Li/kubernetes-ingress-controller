package server

import (
	"context"
	"sync"
)

type IEvent interface {
	Set(ctx context.Context)
	Wait(ctx context.Context)
}

type event struct {
	once sync.Once
	ch   chan struct{}
}

func NewEvent() IEvent {
	return &event{
		ch: make(chan struct{}),
	}
}

func (impl *event) Set(ctx context.Context) {
	impl.once.Do(func() {
		close(impl.ch)
	})
}

func (impl *event) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-impl.ch:
	}
}
