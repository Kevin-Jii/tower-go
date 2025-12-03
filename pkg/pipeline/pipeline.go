package pipeline

// Stage 管道阶段
type Stage[T any] func(input T) (T, error)

// Pipeline 管道
type Pipeline[T any] struct {
	stages []Stage[T]
}

// New 创建管道
func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{}
}

// Add 添加阶段
func (p *Pipeline[T]) Add(stage Stage[T]) *Pipeline[T] {
	p.stages = append(p.stages, stage)
	return p
}

// Execute 执行管道
func (p *Pipeline[T]) Execute(input T) (T, error) {
	var err error
	result := input

	for _, stage := range p.stages {
		result, err = stage(result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// Filter 过滤器（用于切片）
func Filter[T any](input []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range input {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map 映射器
func Map[T, R any](input []T, mapper func(T) R) []R {
	result := make([]R, len(input))
	for i, item := range input {
		result[i] = mapper(item)
	}
	return result
}

// Reduce 归约器
func Reduce[T, R any](input []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, item := range input {
		result = reducer(result, item)
	}
	return result
}

// ChannelPipeline 基于 channel 的管道（并发处理）
type ChannelPipeline[T any] struct {
	input  chan T
	output chan T
}

// NewChannelPipeline 创建 channel 管道
func NewChannelPipeline[T any](bufferSize int) *ChannelPipeline[T] {
	return &ChannelPipeline[T]{
		input:  make(chan T, bufferSize),
		output: make(chan T, bufferSize),
	}
}

// FilterStage 创建过滤阶段
func FilterStage[T any](input <-chan T, predicate func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for item := range input {
			if predicate(item) {
				out <- item
			}
		}
	}()
	return out
}

// MapStage 创建映射阶段
func MapStage[T, R any](input <-chan T, mapper func(T) R) <-chan R {
	out := make(chan R)
	go func() {
		defer close(out)
		for item := range input {
			out <- mapper(item)
		}
	}()
	return out
}
