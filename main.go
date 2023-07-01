package main

import "github.com/mPandaer/go-trace/trace"

func f1() {
	defer trace.Trace("f1")()
	f2()
}

func f2() {
	defer trace.Trace("f2")()
	f3()
}

func f3() {
	defer trace.Trace("f3")()
}

func main() {
	f1()
}
