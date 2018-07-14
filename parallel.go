package parallel

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

var emptyIn = []reflect.Value{}
var emptyContext = context.Background()

//ForLoop type is used in the For function
type ForLoop func(i int)

// panic will not end the program.
// recommend that use the arg ForLoop parameter than using this
func defaultRecover() {
	recover()
}

//For function repeats in parallel, starting with begin and ending with end.
//Internally, it call the ForLoop function each loop
func For(begin int, end int, f ForLoop) {
	ForWithContext(emptyContext, begin, end, f)
}

//ForWithContext function repeats in parallel, starting with begin and ending with end.
//Internally, it call the ForLoop function each loop
func ForWithContext(c context.Context, begin int, end int, f ForLoop) {
	length := end - begin

	if length > 0 {
		ctx, cacnel := context.WithCancel(c)
		go doLoop(cacnel, begin, end, f)
		<-ctx.Done()
	}
}

//doLoop calls the function received as argument in [For]
func doLoop(ctxCancel context.CancelFunc, begin int, end int, f ForLoop) {

	wg := sync.WaitGroup{}
	wg.Add(end - begin)

	for i := begin; i < end; i++ {
		go func(it int) {
			defer wg.Done()
			defer defaultRecover()

			//function call
			f(it)
		}(i)
	}

	wg.Wait()
	ctxCancel()
}

//ForEachSlice loops the slice in parallel
//If put multiple options, only the first one is valid.
//slice: slice, array
//f: any function
//
// s := []int{1,2,3,4,5}
// parallel.ForEachSlice(s, func(i int, e int) {
// 		fmt.Println(i, e)
// })
func ForEachSlice(slice interface{}, f interface{}) {
	ForEachSliceWithContext(emptyContext, slice, f)
}

//ForEachSliceWithContext loops the slice in parallel
//If put multiple options, only the first one is valid.
//slice: slice, array
//f: any function
//
// s := []int{1,2,3,4,5}
// parallel.ForEachSlice(s, func(i int, e int) {
// 		fmt.Println(i, e)
// })
func ForEachSliceWithContext(ctx context.Context, slice interface{}, f interface{}) {
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

		ForWithContext(ctx, 0, reflectionSlice.Len(), func(i int) {
			reflectionFunc.Call([]reflect.Value{reflect.ValueOf(i), reflectionSlice.Index(i)})
		})
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

		ForWithContext(ctx, 0, reflectionSlice.Len(), func(i int) {
			reflectionFunc.Call([]reflect.Value{reflect.ValueOf(i)})
		})
	} else if funcArgc == 0 {
		/**
		* for _ := range slice {
		*	f()
		* }
		**/
		ForWithContext(ctx, 0, reflectionSlice.Len(), func(_ int) {
			reflectionFunc.Call(emptyIn)
		})
	}
}

//ForEachMap loops the Map in parallel
//If put multiple options, only the first one is valid.
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
func ForEachMap(m interface{}, f interface{}) {
	ForEachMapWithContext(emptyContext, m, f)
}

//ForEachMapWithContext loops the Map in parallel
//If put multiple options, only the first one is valid.
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
func ForEachMapWithContext(ctx context.Context, m interface{}, f interface{}) {
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
		ForWithContext(ctx, 0, len(mapKeys), func(i int) {
			key := mapKeys[i]
			reflectionFunc.Call([]reflect.Value{key, reflectionMap.MapIndex(key)})
		})
	} else if funcArgc == 1 {
		/**
		* for k := range m {
		*	f(k)
		* }
		**/
		if keyType, argType := mapType.Key(), funcType.In(0); !keyType.AssignableTo(argType) {
			panic(fmt.Sprintf("map key: %v but function first arg: %v", keyType, argType))
		}
		ForWithContext(ctx, 0, len(mapKeys), func(i int) {
			reflectionFunc.Call([]reflect.Value{mapKeys[i]})
		})
	} else if funcArgc == 0 {
		/**
		* for _ := range m {
		*	f()
		* }
		**/
		ForWithContext(ctx, 0, len(mapKeys), func(_ int) {
			reflectionFunc.Call(emptyIn)
		})
	}

}

//ForEach loops the collection in parallel
//collection: slice, array, map
//If put multiple options, only the first one is valid.
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
func ForEach(collection interface{}, f interface{}) {
	ForEachWithContext(emptyContext, collection, f)
}

//ForEachWithContext loops the collection in parallel
//collection: slice, array, map
//If put multiple options, only the first one is valid.
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
func ForEachWithContext(ctx context.Context, collection interface{}, f interface{}) {
	collectionKind := reflect.TypeOf(collection).Kind()

	switch collectionKind {
	case reflect.Slice, reflect.Array:
		ForEachSliceWithContext(ctx, collection, f)
	case reflect.Map:
		ForEachMapWithContext(ctx, collection, f)
	}
}

//TaskFunc functions that are executed in parallel
type TaskFunc func()

//Race functions that are passed as arguments are executed in parallel,
//and when one of them is finished the function is terminated
// other functions do not force shutdown.
func Race(functions ...TaskFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, e := range functions {
		go func(f TaskFunc) {
			defer cancel()
			defer defaultRecover()

			f()
		}(e)
	}
	<-ctx.Done()
}

//All functions are executed in parallel,
//and when all functions are finished, [All] ends
func All(functions ...TaskFunc) {
	For(0, len(functions), func(i int) {
		functions[i]()
	})
}

//AllWithContext functions are executed in parallel,
//and when all functions are finished, [AllWithContext] ends
//or cancel context called
func AllWithContext(ctx context.Context, functions ...TaskFunc) {
	ForWithContext(ctx, 0, len(functions), func(i int) {
		functions[i]()
	})
}
