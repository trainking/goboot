package component

import (
	"context"
	"sync"
	"time"
)

type (
	// IUpdate 一个Game Update接口定义
	IUpdate interface {

		// Awake 唤醒
		Awake() error

		// Update 执行update
		Update(int32)

		// Destroy 销毁时动作
		Destroy()
	}

	// FixedUpdate 是一个专门模拟物理动作的实现
	FixedUpdate struct {
		ctx        context.Context
		iUpdate    IUpdate   // Update的实现
		deltaTime  float64   // 时间间隔
		lastTime   time.Time // 最后执行时间
		frameIndex int32     // 帧号
		endIndex   int32     // 最后执行帧号，最后一帧是 endIdex -1
		isRunning  bool      // 是否正在运行

		closeChan chan struct{}
		closeOnce sync.Once
	}
)

// NewFixedUpdate 创建一个FixedUpdate
func NewFixedUpdate(ctx context.Context, u IUpdate, deltaTime float64, endIndex int32) *FixedUpdate {
	return &FixedUpdate{
		ctx:       ctx,
		iUpdate:   u,
		deltaTime: deltaTime,
		closeChan: make(chan struct{}),
		endIndex:  endIndex,
	}
}

// Start 开始
func (f *FixedUpdate) Start() error {
	if err := f.iUpdate.Awake(); err != nil {
		return err
	}
	f.lastTime = time.Now()

	go f.run()
	return nil
}

// IsRunning 判断update是否正在运行
func (f *FixedUpdate) IsRunning() bool {
	return f.isRunning
}

// run 执行主体
func (f *FixedUpdate) run() {
	f.isRunning = true
	defer func() {
		recover()
		f.Stop()
		f.isRunning = false
	}()

	for {
		select {
		case <-f.closeChan:
			return
		case <-f.ctx.Done():
			return
		default:
		}

		if f.frameIndex == f.endIndex {
			return
		}

		currentTime := time.Now()
		delta := currentTime.Sub(f.lastTime).Seconds()

		if delta >= f.deltaTime {

			// 执行
			f.iUpdate.Update(f.frameIndex)

			f.frameIndex++
			f.lastTime = currentTime
		}

		time.Sleep(1 * time.Millisecond)
	}
}

// Stop 停止
func (f *FixedUpdate) Stop() {
	f.closeOnce.Do(func() {
		close(f.closeChan)

		// 执行IUpdate实现的销毁
		f.iUpdate.Destroy()
	})
}
