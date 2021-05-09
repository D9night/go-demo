package testing

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

var Wait sync.WaitGroup
var Counter int = 0

func AtomicAddTest(t *testing.T) {
	for routine := 1; routine <= 2; routine ++ {
		Wait.Add(1)
		go Routine(routine)

	}
	Wait.Wait()
	fmt.Printf("Final CounterL: %d\n", Counter)

}

func Routine(id int) {
	for count := 0; count < 2; count++ {
		Counter = Counter + 1
		time.Sleep(1 * time.Nanosecond)
	}
	Wait.Done()
}



func GoWorkGroup(t *testing.T) {
	tr := NewTracker()
	// 是否启动goroutine 应该交给调用者
	go tr.Run()
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	time.Sleep(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
	_ = tr.Event(context.Background(), "test4") // close channel后再发 会painc
}

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{ch: make(chan string, 10)}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}
