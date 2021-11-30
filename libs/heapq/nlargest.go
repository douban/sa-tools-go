package heapq

import (
	"container/heap"
)

func Nlargest(topN uint64, items map[string]uint64) map[string]uint64 {
	pq := make(PriorityQueue, len(items))
	i := 0
	for k, v := range items {
		pq[i] = &Item{
			value:    k,
			priority: v,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)
	ret := make(map[string]uint64)
	count := 0
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		ret[item.value] = item.priority
		count += 1
		if count >= int(topN) {
			break
		}
	}
	return ret
}
