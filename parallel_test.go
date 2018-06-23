package parallel_test

import (
	"fmt"
	"go-parallel"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestForDefault(t *testing.T) {
	parallel.For(0, 3000, func(i int) {
		name := fmt.Sprintf("test%d.txt", i)
		f, _ := os.Open(name)

		f.Write([]byte("hello"))
		f.Close()
		time.Sleep(10 * time.Second)

		os.Remove(name)
	})
}

func TestForLoopArg(t *testing.T) {
	s := make([]int, 1024)
	for i := 0; i < 1024; i++ {
		s[i] = i
	}
	parallel.ForEach(s, func(i int, e int) {
		fmt.Println(i, e)
	}, parallel.TaskOptions{
		TaskCount: 1, // <- single loop..
	})
}

func TestForEachSliceDefault(t *testing.T) {
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEach(s, func(i int, e int) {
		fmt.Println(i, e)
	})
}

func TestForEachSliceInterface(t *testing.T) {
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEach(s, func(i interface{}, val interface{}) {
		fmt.Println(i, val)
	})
}

func TestForEachSliceError(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("int slice - foreach string")
		}
	}()
	s := []int{5, 4, 3, 2, 1}
	// t.
	parallel.ForEach(s, func(i int, e string) {
		fmt.Println(i, e)
	})
}

func TestForEachSliceSingle(t *testing.T) {

	s := []int{5, 4, 3, 2, 1}

	parallel.ForEach(s, func(i int) {
		fmt.Println(i)
	})
}

func TestForEachSliceNoArg(t *testing.T) {

	s := []int{5, 4, 3, 2, 1}

	parallel.ForEach(s, func() {
		fmt.Println("?")
	})
}

func TestForEachMapDefault(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEach(a, func(k string, v int) {
		fmt.Println(k, v)
	})
}

func TestForEachMapDefault2(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k string, v int) {
		fmt.Println(k, v)
	})
}

func TestForEachMapInterface(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k string, v interface{}) {
		fmt.Println(k, v)
	})
}

func TestForEachMapSingle(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k string) {
		fmt.Println(k)
	})
}

func TestForEachMapNoArg(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func() {
		fmt.Println("?")
	})
}

func TestForEachArray(t *testing.T) {

	var a [2048]int
	parallel.ForEach(a, func(i int) {
		a[i] = rand.Int()
	})

	parallel.ForEach(a, func(_ int, e int) {
		fmt.Println(e)
	})
}
