package errgroup

import (
	"context"
	"runtime"
	"sync"

	utilsLogx "github.com/suhanyujie/go_utils/logx"
)

type Group struct {
	err     error
	wg      sync.WaitGroup
	errOnce sync.Once

	workerOnce sync.Once
	ch         chan func(ctx context.Context) error
	chs        []func(ctx context.Context) error

	ctx    context.Context
	cancel func()

	logger *utilsLogx.Logx
}

func WithContext(ctx context.Context, logger *utilsLogx.Logx) *Group {
	return &Group{ctx: ctx, logger: logger}
}

func WithCancel(ctx context.Context, logger *utilsLogx.Logx) *Group {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{ctx: ctx, cancel: cancel, logger: logger}
}

func (g *Group) do(f func(ctx context.Context) error) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			g.logger.Errorf("panic in error group proc, err: %v, stack: %s", r, buf)
		}
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
	err = f(ctx)
}

func (g *Group) GOMAXPROCS(n int) {
	if n <= 0 {
		panic("errgroup: GOMAXPROCS must great than 0")
	}
	g.workerOnce.Do(func() {
		g.ch = make(chan func(context.Context) error, n)
		for i := 0; i < n; i++ {
			go func() {
				for f := range g.ch {
					g.do(f)
				}
			}()
		}
	})
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- f:
		default:
			g.chs = append(g.chs, f)
		}
		return
	}
	go g.do(f)
}

func (g *Group) Wait() error {
	if g.ch != nil {
		for _, f := range g.chs {
			g.ch <- f
		}
	}
	g.wg.Wait()
	if g.ch != nil {
		close(g.ch)
	}
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
