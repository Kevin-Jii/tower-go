package utils

import "sync"

// EventType 事件类型
type EventType string

// Observer 观察者接口
type Observer interface {
	Update(event EventType, data interface{})
}

// EventBus 事件总线
type EventBus struct {
	observers map[EventType][]Observer
	mu        sync.RWMutex
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[EventType][]Observer),
	}
}

// Register 注册观察者
func (e *EventBus) Register(event EventType, observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observers[event] = append(e.observers[event], observer)
}

// Unregister 注销观察者
func (e *EventBus) Unregister(event EventType, observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()

	observers := e.observers[event]
	for i, o := range observers {
		if o == observer {
			e.observers[event] = append(observers[:i], observers[i+1:]...)
			break
		}
	}
}

// Notify 通知所有观察者（同步方式）
func (e *EventBus) Notify(event EventType, data interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	observers := e.observers[event]
	for _, observer := range observers {
		observer.Update(event, data)
	}
}

// NotifyAsync 异步通知所有观察者
func (e *EventBus) NotifyAsync(event EventType, data interface{}) {
	e.mu.RLock()
	observers := e.observers[event]
	e.mu.RUnlock()

	// 启动 goroutine 异步执行
	for _, observer := range observers {
		go observer.Update(event, data)
	}
}

// Event types - 预留事件类型定义
// const (
// 	EventExample EventType = "example.created"
// )

// Global event bus instance
var GlobalEventBus = NewEventBus()
