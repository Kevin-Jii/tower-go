package container

import (
	"sync"
)

// Container 依赖注入容器
type Container struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewContainer 创建容器
func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register 注册服务
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get 获取服务
func (c *Container) Get(name string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.services[name]
}

// MustGet 获取服务，不存在则 panic
func (c *Container) MustGet(name string) interface{} {
	s := c.Get(name)
	if s == nil {
		panic("service not found: " + name)
	}
	return s
}

// Global 全局容器实例
var Global = NewContainer()
