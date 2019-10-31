package utils

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chanxuehong/log"
)

func init() {
	initHook()
}

func initHook() {
	graceExitHook = func() {
		log.Info("grace exit")
	}

	signalWatchedHook = func() {
	}
}

var (
	setupOnce         sync.Once
	globalApp         *App
	graceExitHook     func()
	signalWatchedHook func()
	seq               uint64
)

// App 应用
type App struct {
	ctx context.Context
	sync.Mutex
	children      []*App
	cancel        func()
	wg            sync.WaitGroup
	shutdownHooks []func()
}

// New 创建应用
func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{
		ctx:    ctx,
		cancel: cancel,
	}
	return app
}

// Go 开启应用goroutine
func (app *App) Go(f func(context.Context)) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		f(app.ctx)
	}()
}

// GoWithContext 开启应用goroutine,携带context
func (app *App) GoWithContext(ctx context.Context, f func(context.Context)) {
	app.wg.Add(2)
	wrappedCtx, cancel := context.WithCancel(ctx)
	doneCh := make(chan struct{})
	go func() {
		defer func() {
			app.wg.Done()
		}()
		select {
		case <-ctx.Done():
			cancel()
		case <-app.ctx.Done():
			cancel()
		case <-doneCh:
		}
	}()
	go func() {
		defer func() {
			close(doneCh)
			app.wg.Done()
		}()
		f(wrappedCtx)
	}()
}

// Close 结束这个应用
func (app *App) Close() {
	app.Lock()
	for _, child := range app.children {
		child.Close()
		child.Wait()
	}
	for _, shutdownHook := range app.shutdownHooks {
		shutdownHook()
	}
	app.Unlock()
	app.cancel()
}

// Wait 等待应用goroutine全部退出
func (app *App) Wait() {
	app.wg.Wait()
}

// AddChild 添加子应用
func (app *App) AddChild(child *App) {
	app.Lock()
	app.children = append(app.children)
	app.Unlock()
}

// AddShutdownHook 添加退出钩子
func (app *App) AddShutdownHook(hook func()) {
	app.Lock()
	app.shutdownHooks = append(app.shutdownHooks, hook)
	app.Unlock()
}

// CatchExitSignal 捕获退出信号
func CatchExitSignal(callback func()) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	signalWatchedHook()
	signal := <-signalCh
	log.WithField("signal", signal).
		Info("got-exit-signal")
	callback()
}

// AddShutdownHook 添加退出钩子
func AddShutdownHook(hook func()) {
	setupOnce.Do(setup)
	globalApp.AddShutdownHook(hook)
}

// setup 初始化全局app
func setup() {
	globalApp = New()
	go CatchExitSignal(globalApp.Close)
}

// Go 开启应用goroutine
func Go(f func(context.Context)) {
	setupOnce.Do(setup)
	globalApp.Go(f)
}

// Wait 等待应用goroutine全部退出
func Wait() {
	setupOnce.Do(setup)
	globalApp.Wait()
	graceExitHook()
}

// Close 结束应用
func Close() {
	setupOnce.Do(setup)
	globalApp.Close()
}

// AddChild 添加子应用
func AddChild(child *App) {
	setupOnce.Do(setup)
	globalApp.AddChild(child)
}
