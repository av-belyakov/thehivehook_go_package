package counterelements

import (
	"time"
)

// New новый объект CaseElementCount
// delAfterTime - время, в секундах, по истечении которого объект удаляется
func New(delAfterTime int) *CaseElementCount {
	return &CaseElementCount{
		deletingAfterTime: delAfterTime,
		counters: make(map[string]struct {
			timeExpiry time.Time
			createTime time.Time
			count      int
		})}
}
