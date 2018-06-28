package estimer

import (
	"container/heap"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

const (
	MIN_TIMER_INTERVAL = 1 * time.Millisecond
)

type Timer struct {
	fireTime time.Time
	interval time.Duration
	callback TimerCallback
	repeat   bool
	timerId  uint64
}

func (t *Timer) Cancel() {
	t.callback = nil
}

func (t *Timer) IsActive() bool {
	return t.callback != nil
}

//
// Heap object
//

type _TimerHeap struct {
	timers []*Timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	//log.Println(h.timers[i].fireTime, h.timers[j].fireTime)
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}

	if t1.After(t2) {
		return false
	}
	// t1 == t2, making sure Timer with same deadline is fired according to their add order
	return h.timers[i].timerId < h.timers[j].timerId
}

func (h *_TimerHeap) Swap(i, j int) {
	var tmp *Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

//
// Heap class
//

type HeapTimerQueue struct {
	timerHeap     _TimerHeap
	timerHeapLock sync.Mutex
	timeIdbase    uint64
	timerTable    map[uint64]*Timer
	isExit        bool
	wg            sync.WaitGroup
	ticker        *time.Ticker
}

func NewHeapTimerQueue() *HeapTimerQueue {
	timerQue := new(HeapTimerQueue)
	heap.Init(&timerQue.timerHeap)
	timerQue.timeIdbase = 1
	timerQue.isExit = false
	timerQue.wg.Add(1)
	timerQue.timerTable = make(map[uint64]*Timer)

	go timerQue.TimerLoop(MIN_TIMER_INTERVAL)

	return timerQue
}

// Add a callback which will be called after specified duration
func (this *HeapTimerQueue) NewTimer(delay time.Duration, repeat bool, cb TimerCallback) (uint64, error) {
	t := &Timer{
		fireTime: time.Now().Add(time.Duration(delay)),
		interval: time.Duration(delay),
		callback: cb,
		repeat:   repeat,
	}

	tid := this.timeIdbase
	t.timerId = tid
	this.timeIdbase++

	this.timerHeapLock.Lock()
	heap.Push(&this.timerHeap, t)
	this.timerTable[tid] = t
	this.timerHeapLock.Unlock()

	return tid, nil
}

func (this *HeapTimerQueue) DeleteTimer(tid uint64) {

	if _, ok := this.timerTable[tid]; !ok {
		return
	}

	this.timerTable[tid].Cancel()
}

func (this *HeapTimerQueue) StopTimerQueue() {
	this.isExit = true
	this.wg.Wait()
	fmt.Println("StopTimerQueue Suc!")
}

// Tick once for timers
func (this *HeapTimerQueue) Tick() {
	now := time.Now()
	this.timerHeapLock.Lock()
	for {
		if this.timerHeap.Len() <= 0 {
			break
		}

		nextFireTime := this.timerHeap.timers[0].fireTime

		if nextFireTime.After(now) {
			break
		}

		t := heap.Pop(&this.timerHeap).(*Timer)

		callback := t.callback
		if callback == nil {
			continue
		}

		if !t.repeat {
			t.callback = nil
		}

		this.timerHeapLock.Unlock()
		runCallback(callback)
		this.timerHeapLock.Lock()

		if t.repeat {
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) { // might happen when interval is very small
				t.fireTime = now.Add(t.interval)
			}

			heap.Push(&this.timerHeap, t)
		}
	}

	this.timerHeapLock.Unlock()
}

func (this *HeapTimerQueue) TimerLoop(tickInterval time.Duration) {

	defer this.wg.Done()

	this.ticker = time.NewTicker(time.Millisecond)

	for range this.ticker.C {
		if this.isExit {
			break
		}

		this.Tick()
	}

	this.ticker.Stop()
}

func runCallback(callback TimerCallback) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	callback()
}
