# go-parallel
parallel loop in go 

`go get -u github.com/rudty/go-parallel`

Loop in parallel using goroutine
the function that takes an argument is the same as range

```GO
package main

import (
	"fmt"

	"github.com/rudty/go-parallel"
)

func main() {
	a := []int{1, 2, 3, 4, 5}
	parallel.For(0, len(a), func(i int) {
		fmt.Println(a[i])
	})

	parallel.ForEach(a, func(i, e int) {
		/* prints
		4 5
		0 1
		3 4
		1 2
		2 3
		**/
		fmt.Println(i, e)
	})
}
```

```GO
	//slice
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEach(s, func(i int, e int) {
		fmt.Println(i, e)
	})
```
```GO
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
	//print
	//
	//a 1
	//b b
	//d &{3 4}
	//c {1 2}
```

```GO
	//key value
	s := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEach(s, func(k string, v int) {
		fmt.Println(k, v)
	})
```

```GO
	//key only
	s := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}
	parallel.ForEach(s, func(k string) {
		fmt.Println(k)
	})
```
