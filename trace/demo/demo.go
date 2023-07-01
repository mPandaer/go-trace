package demo

import "trace"

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
