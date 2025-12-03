package plugin

import (
	"fmt"
	"sync"
)

// Plugin 插件接口
type Plugin interface {
	Name() string
	Version() string
	Init() error
	Execute(ctx *Context) error
	Shutdown() error
}

// Context 插件上下文
type Context struct {
	Data   map[string]interface{}
	Errors []error
}

// NewContext 创建上下文
func NewContext() *Context {
	return &Context{
		Data: make(map[string]interface{}),
	}
}

// Set 设置数据
func (c *Context) Set(key string, value interface{}) {
	c.Data[key] = value
}

// Get 获取数据
func (c *Context) Get(key string) interface{} {
	return c.Data[key]
}

// AddError 添加错误
func (c *Context) AddError(err error) {
	c.Errors = append(c.Errors, err)
}

// Kernel 微内核
type Kernel struct {
	plugins  map[string]Plugin
	order    []string
	mu       sync.RWMutex
	started  bool
}

// NewKernel 创建内核
func NewKernel() *Kernel {
	return &Kernel{
		plugins: make(map[string]Plugin),
	}
}

// Register 注册插件
func (k *Kernel) Register(p Plugin) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.started {
		return fmt.Errorf("cannot register plugin after kernel started")
	}

	name := p.Name()
	if _, exists := k.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	k.plugins[name] = p
	k.order = append(k.order, name)
	return nil
}

// Start 启动内核
func (k *Kernel) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	for _, name := range k.order {
		p := k.plugins[name]
		if err := p.Init(); err != nil {
			return fmt.Errorf("failed to init plugin %s: %w", name, err)
		}
	}

	k.started = true
	return nil
}

// Execute 执行所有插件
func (k *Kernel) Execute(ctx *Context) error {
	k.mu.RLock()
	defer k.mu.RUnlock()

	for _, name := range k.order {
		p := k.plugins[name]
		if err := p.Execute(ctx); err != nil {
			ctx.AddError(err)
			// 可以选择继续或中断
		}
	}

	if len(ctx.Errors) > 0 {
		return fmt.Errorf("execution completed with %d errors", len(ctx.Errors))
	}
	return nil
}

// Shutdown 关闭内核
func (k *Kernel) Shutdown() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	// 逆序关闭
	for i := len(k.order) - 1; i >= 0; i-- {
		name := k.order[i]
		p := k.plugins[name]
		if err := p.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown plugin %s: %w", name, err)
		}
	}

	k.started = false
	return nil
}

// GetPlugin 获取插件
func (k *Kernel) GetPlugin(name string) Plugin {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.plugins[name]
}

// ListPlugins 列出所有插件
func (k *Kernel) ListPlugins() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return append([]string{}, k.order...)
}
