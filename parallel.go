package parallel

import (
	"reflect"
	"sync"
)

//TaskOptions options for parallel processing
type TaskOptions struct {

	//TaskCount Sets the goroutine count
	//It only works if it is greater than 0
	TaskCount int

	//PanicHandle call if panic
	PanicHandle func(err interface{})
}

func mixinOptions(opt []TaskOptions) (o TaskOptions) {
	for _, e := range opt {
		if e.TaskCount > 0 {
			o.TaskCount = e.TaskCount
		}
		if e.PanicHandle != nil {
			o.PanicHandle = e.PanicHandle
		}
	}
	return
}

type taskFunc func()

//makeExecutor create worker goroutine
func makeExecutor(c <-chan taskFunc, count int) {
	for i := 0; i < count; i++ {
		go func() {
			for t := range c {
				t()
			}
		}()
	}
}

//ForLoop type is used in the For function
type ForLoop func(i int)

//For function repeats in parallel, starting with begin and ending with end.
//Internally, it call the ForLoop function each loop
func For(begin int, end int, f ForLoop, opt ...TaskOptions) {
	length := end - begin
	if length > 0 {
		option := mixinOptions(opt)

		var executorChan = make(chan taskFunc)
		defer close(executorChan)

		var taskCount int
		if option.TaskCount > 0 {
			taskCount = option.TaskCount
		} else {
			taskCount = length
		}

		makeExecutor(executorChan, taskCount)

		wg := sync.WaitGroup{}
		wg.Add(end - begin)

		var lastError interface{}

		for i := begin; i < end; i++ {
			it := i
			executorChan <- func() {
				defer wg.Done()
				defer func() {
					e := recover()
					if e != nil {
						if option.PanicHandle != nil {
							option.PanicHandle(e)
						} else {
							lastError = e
						}
					}
				}()

				//call function
				f(it)
			}
		}
		wg.Wait()

		if lastError != nil {
			panic(lastError)
		}
	}
}

//ForEachSlice loops the slice in parallel
//slice: slice, array
//f: any function
//
// s := []int{1,2,3,4,5}
// parallel.ForEachSlice(s, func(i int, e int) {
// 		fmt.Println(i, e)
// })
func ForEachSlice(slice interface{}, f interface{}, opt ...TaskOptions) {
	funcType := reflect.TypeOf(f)
	funcArgc := funcType.NumIn()

	reflectionSlice := reflect.ValueOf(slice)
	reflectionFunc := reflect.ValueOf(f)

	if funcArgc == 2 {
		/**
		* for i, e := range slice {
		*	f(i, e)
		* }
		**/
		For(0, reflectionSlice.Len(), func(i int) {
			reflectionFunc.Call([]reflect.Value{reflect.ValueOf(i), reflectionSlice.Index(i)})
		}, opt...)
	} else if funcArgc == 1 {
		/**
		* for i := range slice {
		*	f(i)
		* }
		**/
		For(0, reflectionSlice.Len(), func(i int) {
			reflectionFunc.Call([]reflect.Value{reflect.ValueOf(i)})
		}, opt...)
	} else if funcArgc == 0 {
		/**
		* for _ := range slice {
		*	f()
		* }
		**/
		For(0, reflectionSlice.Len(), func(_ int) {
			reflectionFunc.Call([]reflect.Value{})
		}, opt...)
	}
}

//ForEachMap loops the Map in parallel
//m: map
//f: any function
// a := map[string]int{
// 	"a": 1,
// 	"b": 2,
// 	"c": 3,
// 	"d": 4,
// 	"e": 5,
// }
// parallel.ForEachMap(a, func(k string, v int) {
// 		fmt.Println(k, v)
// })
func ForEachMap(m interface{}, f interface{}, opt ...TaskOptions) {
	funcType := reflect.TypeOf(f)
	funcArgc := funcType.NumIn()

	reflectionMap := reflect.ValueOf(m)
	reflectionFunc := reflect.ValueOf(f)

	mapKeys := reflectionMap.MapKeys()
	if funcArgc == 2 {
		/**
		* for k, v := range m {
		*	f(k, v)
		* }
		**/
		For(0, len(mapKeys), func(i int) {
			key := mapKeys[i]
			reflectionFunc.Call([]reflect.Value{key, reflectionMap.MapIndex(key)})
		}, opt...)
	} else if funcArgc == 1 {
		/**
		* for k := range m {
		*	f(k)
		* }
		**/
		For(0, len(mapKeys), func(i int) {
			reflectionFunc.Call([]reflect.Value{mapKeys[i]})
		}, opt...)
	} else if funcArgc == 0 {
		/**
		* for _ := range m {
		*	f()
		* }
		**/
		For(0, len(mapKeys), func(_ int) {
			reflectionFunc.Call([]reflect.Value{})
		}, opt...)
	}

}

//ForEach loops the collection in parallel
//collection: slice, array, map
//f: any function
//
// ex1)
// s := []int{1,2,3,4,5}
// parallel.ForEach(s, func(i int, e int) {
// 		fmt.Println(i, e)
// })
//
// ex2)
// a := map[string]int{
// 	"a": 1,
// 	"b": 2,
// 	"c": 3,
// 	"d": 4,
// 	"e": 5,
// }
// parallel.ForEach(a, func(k string, v int) {
// 		fmt.Println(k, v)
// })
func ForEach(collection interface{}, f interface{}, opt ...TaskOptions) {
	collectionKind := reflect.TypeOf(collection).Kind()

	switch collectionKind {
	case reflect.Slice, reflect.Array:
		ForEachSlice(collection, f, opt...)
	case reflect.Map:
		ForEachMap(collection, f, opt...)
	}

}
