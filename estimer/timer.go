package estimer

type TimerQueue interface{}
type TimerCallback func()

type TimerInterface interface {
	NewTimer(delay uint32, repeat bool, cb TimerCallback) (uint64, error)
	DeleteTimer(tid uint64) error
	StopTimerQueue() error
}
