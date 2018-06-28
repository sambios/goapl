package estimer

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestCallback(t *testing.T) {
	timer := NewHeapTimerQueue()
	INTERVAL := 100 * time.Millisecond
	for i := 0; i < 10; i++ {
		x := false
		timer.NewTimer(INTERVAL, false, func() {
			fmt.Println("callback!")
			x = true
		})

		time.Sleep(INTERVAL * 2)
		if !x {
			t.Fatalf("x should be true, but it's false")
		}
	}

	// stop timer
	timer.StopTimerQueue()

}

func TestTimer(t *testing.T) {
	timer := NewHeapTimerQueue()
	INTERVAL := 100 * time.Millisecond
	x := 0
	px := x
	now := time.Now()
	nextTime := now.Add(INTERVAL)
	fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)

	timer.NewTimer(INTERVAL, true, func() {
		x += 1
		fmt.Printf("timer %s x %v px %v\n", time.Now(), x, px)
	})

	//time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		time.Sleep(nextTime.Add(INTERVAL / 2).Sub(time.Now()))
		fmt.Printf("Check x %v px %v @ %s\n", x, px, time.Now())
		if x != px+1 {
			t.Fatalf("x should be %d, but it's %d", px+1, x)
		}
		px = x
		nextTime = nextTime.Add(INTERVAL)
		fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)
	}

	timer.StopTimerQueue()
}

func TestCallbackSeq(t *testing.T) {
	timer := NewHeapTimerQueue()
	a := 0
	d := time.Second

	for i := 0; i < 100; i++ {
		i := i
		timer.NewTimer(d, false, func() {
			if a != i {
				t.Error(i, a)
			}

			a += 1
		})
	}
	time.Sleep(d + time.Second*1)

	timer.StopTimerQueue()
}

func TestCancelCallback(t *testing.T) {
	timer := NewHeapTimerQueue()
	INTERVAL := 20 * time.Millisecond
	x := 0

	tid, err := timer.NewTimer(INTERVAL, false, func() {
		x = 1
	})

	if err != nil {
		t.Error("NewTimer failed")
	}

	timer.DeleteTimer(tid)

	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0, but is %v", x)
	}

	// stop timer
	timer.StopTimerQueue()
}

func TestCancelTimer(t *testing.T) {
	timer := NewHeapTimerQueue()
	INTERVAL := 20 * time.Millisecond
	x := 0
	tid, err := timer.NewTimer(INTERVAL, false, func() {
		x += 1
	})

	if err == nil {
		timer.DeleteTimer(tid)
	}

	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0, but is %v", x)
	}

	// stop timer
	timer.StopTimerQueue()
}

func NoTestTimerPerformance(t *testing.T) {
	timer := NewHeapTimerQueue()
	f, err := os.Create("TestTimerPerformance.cpuprof")
	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)
	duration := 10 * time.Second

	for i := 0; i < 400000; i++ {
		if rand.Float32() < 0.5 {
			d := time.Duration(rand.Int63n(int64(duration)))
			timer.NewTimer(d, false, func() {})
		} else {
			d := time.Duration(rand.Int63n(int64(time.Second)))
			timer.NewTimer(d, true, func() {})
		}
	}

	log.Println("Waiting for", duration, "...")
	time.Sleep(duration)
	pprof.StopCPUProfile()
	timer.StopTimerQueue()
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Millisecond * 1500)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
