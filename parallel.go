package parallel

import (
	"context"
	"fmt"
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

//TaskFunc functions that are executed in parallel
type TaskFunc func()

//ForLoop type is used in the For function
type ForLoop func(i int)

//For function repeats in parallel, starting with begin and ending with end.
//Internally, it call the ForLoop function each loop
func For(begin int, end int, f ForLoop, opt ...TaskOptions) {
	length := end - begin
	if length > 0 {
		if len(opt) == 0 {
			forLoopWithoutOption(begin, end, f)
		} else {
			forLoopWithOption(begin, end, f, mixinOptions(opt))
		}
	}
}

//makeExecutor create worker goroutine
func makeExecutor(c <-chan TaskFunc, count int) {
	for i := 0; i < count; i++ {
		go func() {
			for t := range c {
				t()
			}
		}()
	}
}

func forLoopWithOption(begin int, end int, f ForLoop, opt TaskOptions) {

	var executorChan = make(chan TaskFunc)
	defer close(executorChan)

	var taskCount int
	if opt.TaskCount > 0 {
		taskCount = opt.TaskCount
	} else {
		taskCount = end - begin
	}

	makeExecutor(executorChan, taskCount)

	wg := sync.WaitGroup{}
	wg.Add(end - begin)

	for i := begin; i < end; i++ {
		it := i
		executorChan <- func() {
			defer wg.Done()
			defer func() {
				e := recover()
				if e != nil {
					if opt.PanicHandle != nil {
						opt.PanicHandle(e)
					}
				}
			}()

			//call function
			f(it)
		}
	}
	wg.Wait()
}

func forLoopWithoutOption(begin int, end int, f ForLoop) {

	wg := sync.WaitGroup{}
	wg.Add(end - begin)

	for i := begin; i < end; i++ {
		go func(idx int) {
			defer wg.Done()
			defer func() {
				//recevie uncaught panic
				recover()
			}()

			//call function
			f(idx)
		}(i)
	}

	wg.Wait()
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
	reflectionSlice := reflect.ValueOf(slice)
	reflectionFunc := reflect.ValueOf(f)

	if reflectionSlice.Len() == 0 {
		return
	}

	funcType := reflect.TypeOf(f)
	funcArgc := funcType.NumIn()

	sliceType := reflect.TypeOf(slice)

	if funcArgc == 2 {
		/**
		* for i, e := range slice {
		*	f(i, e)
		* }
		**/

		if !reflect.TypeOf(0).AssignableTo(funcType.In(0)) {
			//reflect.TypeOf(0) = int type
			panic("first argument is not an int")
		}

		if elemType, argType := sliceType.Elem(), funcType.In(1); !elemType.AssignableTo(argType) {
			panic(fmt.Sprintf("slice value type: %v but func second arg type: %v", elemType, argType))
		}

		For(0, reflectionSlice.Len(), func(i int) {
			reflectionFunc.Call([]reflect.Value{reflect.ValueOf(i), reflectionSlice.Index(i)})
		}, opt...)
	} else if funcArgc == 1 {
		/**
		* for i := range slice {
		*	f(i)
		* }
		**/

		if !reflect.TypeOf(0).AssignableTo(funcType.In(0)) {
			//reflect.TypeOf(0) = int type
			panic("first argument is not an int")
		}

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
	reflectionMap := reflect.ValueOf(m)

	if reflectionMap.Len() == 0 {
		return
	}

	reflectionFunc := reflect.ValueOf(f)

	funcType := reflect.TypeOf(f)
	funcArgc := funcType.NumIn()

	mapType := reflectionMap.Type()
	mapKeys := reflectionMap.MapKeys()
	if funcArgc == 2 {
		/**
		* for k, v := range m {
		*	f(k, v)
		* }
		**/
		if keyType, argType := mapType.Key(), funcType.In(0); !keyType.AssignableTo(argType) {
			panic(fmt.Sprintf("map keyType: %v but func first argType: %v", keyType, argType))
		}
		if valType, argType := mapType.Elem(), funcType.In(1); !valType.AssignableTo(argType) {
			panic(fmt.Sprintf("map valueType: %v but func second argType: %v", valType, argType))
		}
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
		if keyType, argType := mapType.Key(), funcType.In(0); !keyType.AssignableTo(argType) {
			panic(fmt.Sprintf("map key: %v but function first arg: %v", keyType, argType))
		}
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

//Race functions that are passed as arguments are executed in parallel,
//and when one of them is finished the function is terminated
// other functions do not force shutdown.
func Race(functions ...TaskFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, e := range functions {
		go func(f TaskFunc) {
			defer cancel()
			defer func() {
				//recevie uncaught panic
				recover()
			}()

			f()
		}(e)
	}
	<-ctx.Done()
}
