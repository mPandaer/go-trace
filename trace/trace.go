package trace

import "fmt"

func Trace(name string) func() {
	fmt.Printf("enter: %s\n", name)
	return func() {
		fmt.Printf("exit: %s\n", name)
	}
}
