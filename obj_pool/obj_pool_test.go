package obj_pool_test

import (
	"runtime"
	"testing"

	baseObjPool "github.com/lkjx82/basego/obj_pool"

	"fmt"
	"math/rand"
	"time"
)

var makes int
var frees int

func makeBuffer() interface{} {
	makes += 1
	return make([]byte, 5000000+5000000)
}

// -----------------------------------------------------------------------------
//
func TestObjectPoolFixed(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	objPool := baseObjPool.NewObjPoolFixed(10, makeBuffer, true)

	testCache := make([][]byte, 20)

	var m runtime.MemStats
	for i := 0; i < 1000; i++ {
		b := objPool.Alloc().([]byte)

		i := rand.Intn(len(testCache))
		if testCache[i] != nil {
			objPool.Free(testCache[i])
		}

		testCache[i] = b

		time.Sleep(time.Millisecond * time.Duration(50+rand.Int63n(50)))

		bytes := 0
		for i := 0; i < len(testCache); i++ {
			if testCache[i] != nil {
				bytes += len(testCache[i])
			}
		}

		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes, frees)
	}

	return
}

// -----------------------------------------------------------------------------
//
func dTestObjectPool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	objPool := baseObjPool.NewObjPool(makeBuffer)

	testCache := make([][]byte, 20)

	var m runtime.MemStats
	for i := 0; i < 1000; i++ {
		b := objPool.Alloc().([]byte)

		i := rand.Intn(len(testCache))
		if testCache[i] != nil {
			objPool.Free(testCache[i])
		}

		testCache[i] = b

		time.Sleep(time.Millisecond * time.Duration(50+rand.Int63n(50)))

		bytes := 0
		for i := 0; i < len(testCache); i++ {
			if testCache[i] != nil {
				bytes += len(testCache[i])
			}
		}

		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes, frees)
	}

	return
}

// -----------------------------------------------------------------------------
//func BenchmarkObjectPool(b *testing.B) {
//	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

//	objPool := base.NewObjPool(makeBuffer)

//	testCache := make([][]byte, 20)

//	var m runtime.MemStats
//	for i := 0; i < b.N; i++ {
//		b := objPool.Alloc().([]byte)

//		i := rand.Intn(len(testCache))
//		if testCache[i] != nil {
//			objPool.Free(testCache[i])
//		}

//		testCache[i] = b

//		time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)))

//		bytes := 0
//		for i := 0; i < len(testCache); i++ {
//			if testCache[i] != nil {
//				bytes += len(testCache[i])
//			}
//		}

//		runtime.ReadMemStats(&m)
//		fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
//			m.HeapIdle, m.HeapReleased, makes, frees)
//	}

//	return
//}
