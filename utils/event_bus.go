package utils

import "sync"

// EventType 事件类型
type EventType string

// 预定义事件类型
const (
	EventOrderCreated   EventType = "order.created"
	EventOrderConfirmed EventType = "order.confirmed"
	EventOrderCompleted EventType = "order.completed"
	EventOrderCancelled EventType = "order.cancelled"
	EventSupplierBound  EventType = "supplier.bound"
)

// EventHandler 事件处理函数类型
type EventHandler func(data interface{})

// Observer 观察者接口（保留向后兼容）
type Observer interface {
	Update(event EventType, data interface{})
}

// EventBus 事件总线
type EventBus struct {
	observers map[EventType][]Observer
	handlers  map[EventType][]EventHandler
	mu        sync.RWMutex
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[EventType][]Observer),
		handlers:  make(map[EventType][]EventHandler),
	}
}

// Register 注册观察者
func (e *EventBus) Register(event EventType, observer Observer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observers[event] = append(e.observers[event], observer)
}

// Subscribe 订阅事件（函数式，更简洁）
func (e *EventBus) Subscribe(event EventType, handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
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

// Publish 发布事件（异步）
func (e *EventBus) Publish(event EventType, data interface{}) {
	e.mu.RLock()
	observers := e.observers[event]
	handlers := e.handlers[event]
	e.mu.RUnlock()

	for _, observer := range observers {
		go observer.Update(event, data)
	}
	for _, handler := range handlers {
		go handler(data)
	}
}

// PublishSync 发布事件（同步）
func (e *EventBus) PublishSync(event EventType, data interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, observer := range e.observers[event] {
		observer.Update(event, data)
	}
	for _, handler := range e.handlers[event] {
		handler(data)
	}
}

// Notify 通知所有观察者（同步方式，保留向后兼容）
func (e *EventBus) Notify(event EventType, data interface{}) {
	e.PublishSync(event, data)
}

// NotifyAsync 异步通知所有观察者（保留向后兼容）
func (e *EventBus) NotifyAsync(event EventType, data interface{}) {
	e.Publish(event, data)
}

// Global event bus instance
var GlobalEventBus = NewEventBus()
