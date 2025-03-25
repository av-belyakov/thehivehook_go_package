package webhookserver

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	sfunc "github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

var (
	once     sync.Once
	memCache *MemoryCache = nil
)

type MemoryCache struct {
	Alloc                  uint64
	HeapSys                uint64
	HeapAlloc              uint64
	TotalAlloc             uint64
	HeapObjects            uint64
	NumberLiveObjects      uint64
	CountMemoryReturned    uint64
	GarbagecollectorMemory uint64
}

func NewMemoryCache() *MemoryCache {
	once.Do(func() {
		memCache = new(MemoryCache)
	})

	return memCache
}

// printMemStats вывод информации по потребляемой памяти
func printMemStats() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memCache = NewMemoryCache()

	alloc := memStats.Alloc
	gcSys := memStats.GCSys
	heapSys := memStats.HeapSys
	heapAlloc := memStats.HeapAlloc
	numLiveObj := memStats.Mallocs - memStats.Frees
	returnedOS := memStats.HeapIdle - memStats.HeapReleased
	totalAlloc := memStats.TotalAlloc
	heapObjects := memStats.HeapObjects

	str := strings.Builder{}

	str.WriteString(fmt.Sprintf("Allocated Memory: %v bytes %s\n", alloc, sfunc.GetPointerUpOrDown(memCache.Alloc, alloc)))
	str.WriteString(fmt.Sprintf("Total Allocated Memory: %v bytes %s\n", totalAlloc, sfunc.GetPointerUpOrDown(memCache.TotalAlloc, totalAlloc)))
	str.WriteString(fmt.Sprintf("Heap Alloc Memory: %v bytes %s\n", heapAlloc, sfunc.GetPointerUpOrDown(memCache.HeapAlloc, heapAlloc)))
	str.WriteString(fmt.Sprintf("Heap System Memory: %v bytes %s\n", heapSys, sfunc.GetPointerUpOrDown(memCache.HeapSys, heapSys)))
	str.WriteString(fmt.Sprintf("The number of allocated heap objects: %v bytes %s\n", heapObjects, sfunc.GetPointerUpOrDown(memCache.HeapObjects, heapObjects)))
	str.WriteString(fmt.Sprintf("The number of live objects: %v bytes %s\n", numLiveObj, sfunc.GetPointerUpOrDown(memCache.NumberLiveObjects, numLiveObj)))
	str.WriteString(fmt.Sprintf("Count memory that could be returned to the OS: %v bytes %s\n", returnedOS, sfunc.GetPointerUpOrDown(memCache.CountMemoryReturned, returnedOS)))
	str.WriteString(fmt.Sprintf("Garbage Collector Memory: %v bytes %s\n", gcSys, sfunc.GetPointerUpOrDown(memCache.GarbagecollectorMemory, gcSys)))

	memCache.Alloc = alloc
	memCache.HeapSys = heapSys
	memCache.HeapAlloc = heapAlloc
	memCache.TotalAlloc = totalAlloc
	memCache.HeapObjects = heapObjects
	memCache.NumberLiveObjects = numLiveObj
	memCache.CountMemoryReturned = returnedOS
	memCache.GarbagecollectorMemory = gcSys

	return str.String()
}
