# go-parallel
parallel loop in go 

`go get -u github.com/rudty/go-parallel`

Loop in parallel using goroutine
the function that takes an argument is the same as range

```GO
	//slice
	s := []int{5, 4, 3, 2, 1}
	parallel.ForEach(s, func(i int, e int) {
		fmt.Println(i, e)
	})
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
