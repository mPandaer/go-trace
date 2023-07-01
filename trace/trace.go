package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("产生了一些错误 in Trance()")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	id := curGoroutineId()
	mu.Lock()
	indent := m[id]
	m[id] = indent + 1
	mu.Unlock()
	printTrace(id, name, "->", indent)
	return func() {
		mu.Lock()
		indent = m[id]
		m[id] = indent - 1
		mu.Unlock()
		printTrace(id, name, "<-", indent-1)
	}
}

var goroutineSpace = []byte("goroutine ")

func curGoroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space find in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q %v", b, err))
	}
	return n

}

func printTrace(gid uint64, name string, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "      "
	}
	fmt.Printf("G[%05d]: %s%s%s\n", gid, indents, arrow, name)
}
