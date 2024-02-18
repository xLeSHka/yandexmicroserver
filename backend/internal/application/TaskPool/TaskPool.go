package taskpool

import (
	"fmt"
	"sync"
)

// type Task struct {
// 	Expression string
// }

// func (t *Task) Process() {

// }

// type WorkerPool struct {
// 	Tasks   []Task
// 	goCount int
// 	tasksCh chan Task
// 	wg      sync.WaitGroup
// }

// func (wp *WorkerPool) worker() {
// 	for task := range wp.tasksCh {
// 		task.Process()
// 		wp.wg.Done()
// 	}
// }

// func (wp *WorkerPool) Run() {

// }

// type Task struct {
// 	Address    string
// 	ctx        context.Context
// 	log        *slog.Logger
// 	exprassion string
// 	client     *http.Client
// }

// func (t *Task) Process() (float64, error) {
// 	data := []byte(t.exprassion)
// 	r := bytes.NewReader(data)
// 	resp, err := http.Post(t.Address, "", r)
// 	if err != nil {
// 		t.log.Error("failed post subExprassion to agent: %v", err)
// 		return 0, err
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		t.log.Error("failed read resp body: %v", err)
// 		return 0, err
// 	}
// 	res, err := strconv.ParseFloat(string(body), 256)
// 	if err != nil {
// 		t.log.Error("failed parse resp body from bytes to float: %v", err)
// 		return 0, err
// 	}
// 	return res, nil

// }

// type WorkerPool struct {
// 	Tasks       chan Task
// 	AgentsCount int
// 	wg          sync.WaitGroup
// }

// func (wp *WorkerPool) worker() {
// 	for task := range wp.Tasks {
// 		task.Process()
// 		wp.wg.Done()
// 	}
// }

type PoolTask interface {
	Execute() error
	OnFailure(error)
}
type CalcTask struct {
	Exprassion string
}
type WorkerPool interface {
	Start()
	Stop()
	AddWork(t PoolTask)
}

type MyPool struct {
	tasks      chan PoolTask
	numWorkers int
	wg         *sync.WaitGroup
	mu         sync.RWMutex
	started    bool
	stopped    bool
}

func NewWorkerPool(numWorkers int, channelSize int) (*MyPool, error) {
	if numWorkers <= 0 {
		return nil, fmt.Errorf("incoorect numWorkers")
	}

	if channelSize < 0 {
		return nil, fmt.Errorf("negative channelSize")
	}
	return &MyPool{tasks: make(chan PoolTask, channelSize),
		numWorkers: numWorkers,
		wg:         &sync.WaitGroup{}}, nil
}

func (p *MyPool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.started || p.stopped {
		return
	}
	p.started = true
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.tasks {
				if err := task.Execute(); err != nil {
					task.OnFailure(err)
				}
			}
		}()
	}
}

func (p *MyPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped || !p.started {
		return
	}
	close(p.tasks)
	p.stopped = true
	p.wg.Wait()
}

func (p *MyPool) AddWork(t PoolTask) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.stopped {
		return
	}
	p.tasks <- t
}
