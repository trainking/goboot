package errgroup

import (
	"context"
	"sync"
)

// ForIndexGroup 是一个引用for index执行多个携程的ErrGroup
type ForIndexGroup struct {
	err     error
	wg      sync.WaitGroup
	errOnce sync.Once

	ctx    context.Context
	cancel func()
}

// WithContextForIndexGroup 创建ForIndexGroup带Group
func WithContextForIndexGroup(ctx context.Context) *ForIndexGroup {
	return &ForIndexGroup{ctx: ctx}
}

// WithCancelForIndexGroup 创建ForIndexGroup带Context和Cancel
func WithCancelForIndexGroup(ctx context.Context) *ForIndexGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &ForIndexGroup{ctx: ctx, cancel: cancel}
}

// Go 启动携程
func (g *ForIndexGroup) Go(f func(ctx context.Context, i int) error, i int) {
	g.wg.Add(1)
	go g.do(f, i)
}

// do 执行携程
func (g *ForIndexGroup) do(f func(ctx context.Context, i int) error, i int) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	var err error
	defer func() {
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f(ctx, i)
}

// Wait 阻塞得到结果
func (g *ForIndexGroup) Wait() error {
	g.wg.Wait()
	return g.err
}
