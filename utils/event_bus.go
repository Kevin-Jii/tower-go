package utils

import (
	"sync"
)

// Event 事件接口
type Event interface {
	Name() string
}

// EventHandler 事件处理器类型
type EventHandler func(event Event) error

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

var (
	globalEventBus     *EventBus
	globalEventBusOnce sync.Once
)

// GetEventBus 获取全局事件总线实例
func GetEventBus() *EventBus {
	globalEventBusOnce.Do(func() {
		globalEventBus = NewEventBus()
	})
	return globalEventBus
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, exists := eb.handlers[eventName]; !exists {
		eb.handlers[eventName] = make([]EventHandler, 0)
	}
	eb.handlers[eventName] = append(eb.handlers[eventName], handler)
}

// Publish 发布事件（同步）
func (eb *EventBus) Publish(event Event) []error {
	eb.mu.RLock()
	handlers, exists := eb.handlers[event.Name()]
	eb.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return nil
	}

	var errors []error
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// PublishAsync 发布事件（异步）
func (eb *EventBus) PublishAsync(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.handlers[event.Name()]
	eb.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			_ = h(event)
		}(handler)
	}
}

// Unsubscribe 取消订阅（移除所有处理器）
func (eb *EventBus) Unsubscribe(eventName string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.handlers, eventName)
}

// HasSubscribers 检查事件是否有订阅者
func (eb *EventBus) HasSubscribers(eventName string) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	handlers, exists := eb.handlers[eventName]
	return exists && len(handlers) > 0
}
