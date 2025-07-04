package scheduler

type Item struct {
	Value    string // 실제 저장할 값 (예: 작업 이름)
	Priority int    // 우선순위 값 (클수록 우선)
	Index    int    // heap 내 인덱스 (heap.Fix 등에서 필요)
}

type PriorityQueue []*Item

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].Priority > (*pq)[j].Priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // 안전한 상태로 초기화
	*pq = old[0 : n-1]
	return item
}
