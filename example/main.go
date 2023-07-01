package main

import (
	"sync"
	"trace"
)

func f1() {
	defer trace.Trace()()
	f2()
}

func f2() {
	defer trace.Trace()()
	f3()
}

func f3() {
	defer trace.Trace()()
}

func g1() {
	defer trace.Trace()()
	g2()
}

func g2() {
	defer trace.Trace()()
	g3()
}

func g3() {
	defer trace.Trace()()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		f1()
		wg.Done()
	}()

	g1()
	wg.Wait()
}
