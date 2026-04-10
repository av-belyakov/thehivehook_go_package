package counterelements

import (
	"sync"
	"time"
)

// CaseElementCount объект для подсчета количества элементов
type CaseElementCount struct {
	mtx      sync.RWMutex
	counters map[string]struct {
		timeExpiry time.Time
		createTime time.Time
		count      int
	}
	deletingAfterTime int
}
