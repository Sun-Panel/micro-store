package event

import "sync"

// Handler 处理单个事件载荷
type Handler func(payload interface{})

type envelope struct {
	name    string
	payload interface{}
}

// Dispatcher 极简进程内事件分发器：buffered channel + 固定 worker 池。
// 适用于单实例部署；多实例需替换为外部 broker（Redis/NSQ）以避免重复消费。
type Dispatcher struct {
	ch       chan envelope
	handlers map[string][]Handler
	mu       sync.RWMutex
	wg       sync.WaitGroup
	workers  int
	// OnDrop 队列满丢弃时的回调（可选），用于告警/排查
	OnDrop func(name string)
}

// New 创建分发器，workers 个消费协程（<=0 时默认 4）
func New(workers int) *Dispatcher {
	if workers <= 0 {
		workers = 4
	}
	return &Dispatcher{
		ch:       make(chan envelope, 1024),
		handlers: make(map[string][]Handler),
		workers:  workers,
	}
}

// Register 注册某事件的 handler（同一事件可注册多个）
func (d *Dispatcher) Register(name string, h Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[name] = append(d.handlers[name], h)
}

// Publish 非阻塞发布：入队即返回；队列满则触发 OnDrop 并丢弃（不让调用方阻塞）
func (d *Dispatcher) Publish(name string, payload interface{}) {
	select {
	case d.ch <- envelope{name: name, payload: payload}:
	default:
		if d.OnDrop != nil {
			d.OnDrop(name)
		}
	}
}

// Start 启动 worker 池，开始消费事件
func (d *Dispatcher) Start() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.loop()
	}
}

func (d *Dispatcher) loop() {
	defer d.wg.Done()
	for e := range d.ch {
		d.dispatch(e)
	}
}

// dispatch 逐个执行 handler，单个 handler panic 不影响其余
func (d *Dispatcher) dispatch(e envelope) {
	d.mu.RLock()
	hs := make([]Handler, len(d.handlers[e.name]))
	copy(hs, d.handlers[e.name])
	d.mu.RUnlock()
	for _, h := range hs {
		func(h Handler) {
			defer func() { _ = recover() }()
			h(e.payload)
		}(h)
	}
}

// Stop 优雅关闭：关闭入队通道并等待已入队事件处理完毕
func (d *Dispatcher) Stop() {
	close(d.ch)
	d.wg.Wait()
}

// ---- 默认实例与包级便捷函数 ----

// Default 进程内默认分发器
var Default = New(8)

// Register 在默认分发器上注册 handler
func Register(name string, h Handler) { Default.Register(name, h) }

// Publish 通过默认分发器发布事件
func Publish(name string, payload interface{}) { Default.Publish(name, payload) }

// Start 启动默认分发器
func Start() { Default.Start() }

// Stop 优雅关闭默认分发器
func Stop() { Default.Stop() }
