package parallel_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rudty/go-parallel"
)

func TestForDefault(t *testing.T) {
	parallel.For(0, 3000, func(i int) {
		name := fmt.Sprintf("test%d.txt", i)
		f, err := os.Create(name)
		if err != nil {
			t.Error(err)
		}

		f.Write([]byte("hello"))
		f.Close()
		time.Sleep(1 * time.Second)

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
	})
}

func TestForEachSliceDefault(t *testing.T) {
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEach(s, func(i int, e int) {
		fmt.Println(i, e)
	})

}

func TestForEachSliceDefault2(t *testing.T) {
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEachSlice(s, func(i int, e int) {
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

func TestForEachSliceError2(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("int slice - foreach string")
		}
	}()
	s := []int{5, 4, 3, 2, 1}
	// t.
	parallel.ForEach(s, func(i string, e string) {
		fmt.Println(i, e)
	})
}

func TestForEachSliceError3(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("int slice - foreach string")
		}
	}()
	s := []int{5, 4, 3, 2, 1}
	// t.
	parallel.ForEach(s, func(i string) {
		fmt.Println(i)
	})
}

func TestForEachSliceEmpty(t *testing.T) {

	s := []int{}

	parallel.ForEach(s, func(i int) {
		fmt.Println("?")
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

func TestForEachMapKeyError(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("string key func int")
		}
	}()

	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k int) {
		fmt.Println("?", k)
	})
}

func TestForEachMapValueError(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("int val func string")
		}
	}()

	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k string, v string) {
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

func TestForEachMapEmpty(t *testing.T) {
	a := map[string]int{}
	parallel.ForEachMap(a, func() {
		fmt.Println("?")
	})
}

func TestForEachMapBadKey(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("string val func int")
		}
	}()
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(a int) {
		fmt.Println("?")
	})
}

func TestForEachMapBadKey2(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("string val func int")
		}
	}()
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEachMap(a, func(k int, v interface{}) {
		fmt.Println("?")
	})
}

func TestRaceDefault(t *testing.T) {
	parallel.Race(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("1sec")
	},
		func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("500msec")
		},
		func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Println("50msec")
		})
	fmt.Println("end Rance")

}

func TestForEachMapAny(t *testing.T) {
	type foo struct {
		a int
		b int
	}
	a := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": foo{1, 2},
		"d": &foo{3, 4},
	}
	parallel.ForEach(a, func(k string, v interface{}) {
		fmt.Println(k, v)
	})

}

func TestAllDefault(t *testing.T) {
	parallel.All(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("1sec")
	},
		func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("500msec")
		},
		func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Println("50msec")
		})
	fmt.Println("end All")

}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	parallel.ForWithContext(ctx, 0, 100, func(i int) {
		// time.Sleep(1 * time.Second)
	})

	if ctx.Err() != nil {
		t.Error(ctx.Err())
	}
}

func TestContext2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	parallel.ForWithContext(ctx, 0, 100, func(i int) {
		time.Sleep(1 * time.Second)
	})

	if ctx.Err() == nil {
		t.Error("require timeout error")
	}
}

func TestAllWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	parallel.AllWithContext(ctx, func() {

	}, func() {

	})

	if ctx.Err() != nil {
		t.Error(ctx.Err())
	}
}

func TestAllWithContext2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	parallel.AllWithContext(ctx, func() {

	}, func() {
		time.Sleep(2 * time.Second)
	})

	if ctx.Err() == nil {
		t.Error("require timeout error")
	}
}
