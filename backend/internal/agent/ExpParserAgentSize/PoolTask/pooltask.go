package pooltask

// import (
// 	"errors"
// 	"sync"

// 	exp "github.com/xleshka/distributedcalc/backend/internal/application/ExpParser"
// )

// type Task struct {
// 	Err error
// 	Data interface{}
// 	f func(interface{}) error
// }
// func NewTask(f func(interface{})error,data interface{}) *Task {
// 	return &Task{f: f,Data: data}
// }
// func Execute(task *Task) {
// 	task.Err = task.f(task.Data)
// }

// type WorkerPool interface {
// 	/* Start подготавливает пул для обработки задач. */
// 	Start()
// 	/* Stop останавливает обработку в пуле.  */
// 	Stop()
// 	/* AddWork добавляет задачу для обработки пулом.  */
// 	AddWork(t PoolTask)
// }
// type Pool struct {
// 	started  bool
// 	stopped bool
// 	tasks chan PoolTask
// 	numWorkers int
// 	mu       sync.RWMutex
// 	wg *sync.WaitGroup
// }
// func NewWorkerPool(numWorkers int, channelSize int) (*Pool,error) {
// 	if numWorkers <= 0 {
// 		return nil,errors.New("Null workers")
// 	}
// 	if channelSize < 0 {
// 		return nil, errors.New("Negative chan size")
// 	}

// 	return &Pool{
// 		tasks: make(chan PoolTask,channelSize),
// 		numWorkers: numWorkers,
// 		wg: &sync.WaitGroup{},
// 	}, nil
// }
// func (p *Pool)Start() {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	if p.started || p.stopped {
// 		return
// 	}
// 	p.started = true
// 	for i := 0; i < p.numWorkers;i++ {
// 		p.wg.Add(1)
// 		go func() {
// 			defer p.wg.Done()
// 			for _, t := range p.tasks
// 		}
// 	}
// }
