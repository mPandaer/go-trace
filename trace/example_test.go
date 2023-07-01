package trace

import "testing"

func f1() {
	defer Trace()()
	f2()
}

func f2() {
	defer Trace()()
	f3()
}

func f3() {
	defer Trace()()
}

func TestExample(t *testing.T) {
	f1()
}
